package admin

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	gin_zap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/chaosnote/wander/utils"
	"github.com/chaosnote/wander/web/admin/router/version"
)

//-----------------------------------------------

func gen_logger(log_mode int, log_path string) *zap.Logger {
	var logger *zap.Logger
	switch log_mode {
	case 0:
		logger = utils.NewConsoleLogger(1)
	case 1:
		logger = utils.NewFileLogger(log_path, 1)
	default:
		panic(fmt.Errorf("unknow log mode (%v)", log_mode))
	}
	return logger
}

//-----------------------------------------------

type AdminStore interface {
	Start()
	Close()
}

type admin_store struct {
}

func (s *admin_store) Start() {
	var di = utils.GetDI()
	logger := di.MustGet(LOGGER_SYSTEM).(*zap.Logger)

	middleware := NewMiddlewareStore()
	redis_store, _ := redis.NewStore(10, "tcp", "192.168.0.236:6379", "", "", []byte("secret"))

	g := gin.New()
	g.Use(sessions.Sessions("custom_session", redis_store))
	g.Use(gin_zap.RecoveryWithZap(utils.NewConsoleLogger(1), true))
	g.Use(middleware.Logging())

	g.GET("/version", version.Output)

	const addr = ":9090"
	for _, route := range g.Routes() {
		logger.Debug(fmt.Sprintf("http://localhost%s%s", addr, route.Path))
	}

	server := &http.Server{
		Addr:    addr,
		Handler: g,
	}

	go func() {
		e := server.ListenAndServe()
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (s *admin_store) Close() {
	var di = utils.GetDI()

	di.MustGet(LOGGER_SYSTEM).(*zap.Logger).Sync()
	di.MustGet(LOGGER_GIN).(*zap.Logger).Sync()
}

//-----------------------------------------------

func NewAdminStore() AdminStore {
	var s = &admin_store{}

	flag.Parse()

	log_mode := *LOG_MODE
	switch log_mode {
	case 0:
	case 1:
		gin.SetMode(gin.DebugMode)
	case 2:
		gin.SetMode(gin.ReleaseMode)
	default:
		panic(fmt.Errorf("unknow log mode (%v)", log_mode))
	}

	var di = utils.GetDI()
	di.SetShare(LOGGER_SYSTEM, func(args ...interface{}) any {
		return gen_logger(log_mode, "./logs/admin/system")
	})
	di.SetShare(LOGGER_GIN, func(args ...interface{}) any {
		return gen_logger(log_mode, "./logs/admin/gin")
	})
	return s
}
