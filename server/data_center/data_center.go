package datacenter

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/chaosnote/wander/data_center/api"
	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/utils"
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

type DCStore interface {
	Start()
	Close()
}

//-----------------------------------------------

type store struct {
	logger *zap.Logger

	DBStore
	NatsStore
	RedisStore
	PlayerStore
	api.APIStore
}

//-----------------------------------------------

func (s *store) Start() {
	const msg = "Start"
	var e error

	middleware := NewMiddlewareStore()
	// router
	router := mux.NewRouter()
	router.Use(middleware.Logging)

	// 更新值[未處理]
	sub := router.PathPrefix("/api").Subrouter()

	sub = router.PathPrefix("/guest").Subrouter()
	sub.Use(middleware.Guest)
	sub.HandleFunc(`/new`, s.HandleGuestNew).Methods(http.MethodGet)

	sub = router.PathPrefix("/player").Subrouter()
	sub.HandleFunc(`/login`, s.HandlePlayerLogin).Methods(http.MethodPost)
	sub.HandleFunc(`/logout`, s.HandlePlayerLogout).Methods(http.MethodPost)

	e = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, e := route.GetPathTemplate()
		if e != nil {
			return e
		}
		s.logger.Debug(msg, zap.String("path", template))
		return nil
	})
	if e != nil {
		panic(e)
	}

	s.logger.Debug(fmt.Sprintf("addr %s", dc_addr))
	go func() {
		e = http.ListenAndServe(dc_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (ds *store) Close() {
	var di = utils.GetDI()
	di.MustGet(internal.SERVICE_MARIADB).(*sql.DB).Close()
	di.MustGet(internal.SERVICE_NATS).(*nats.Conn).Close()
	di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger).Sync()
}

//-----------------------------------------------

func NewDCStore() DCStore {
	flag.Parse()

	utils.RSAInit("./asset/_rsa/pem.key", 512, true)

	var s = &store{}
	var di = utils.GetDI()
	var e error

	log_mode := *LOG_MODE
	di.SetShare(internal.LOGGER_SYSTEM, func(...interface{}) any {
		return gen_logger(log_mode, "./logs/data_center/system")
	})
	di.Set(internal.LOGGER_API, func(args ...interface{}) any {
		dir_path := args[0].(string)
		if len(dir_path) == 0 {
			panic("unset dir path")
		}
		return gen_logger(log_mode, filepath.Join("./logs/data_center", dir_path))
	})

	di.SetShare(internal.SERVICE_MARIADB, func(...interface{}) any {
		// 例 : "user:password@tcp(ip)?parseTime=true/dbname"
		cmd := fmt.Sprintf(`%s:%s@tcp(%s)/%s?parseTime=true`, db_user, db_pw, db_addr, db_name)
		var db *sql.DB
		db, e = sql.Open("mysql", cmd)
		if e != nil {
			panic(e)
		}
		e = db.Ping()
		if e != nil {
			panic(e)
		}
		db.SetMaxOpenConns(100)                // Limit to N open connections
		db.SetMaxIdleConns(10)                 // Keep up to N idle connections
		db.SetConnMaxLifetime(5 * time.Minute) // Reuse connections for at most N

		return db
	})

	di.SetShare(internal.SERVICE_NATS, func(...interface{}) any {
		var conn *nats.Conn
		conn, e = nats.Connect(fmt.Sprintf("nats://%s", nats_addr))
		if e != nil {
			panic(e)
		}
		return conn
	})

	di.SetShare(internal.SERVICE_REDIS, func(...interface{}) any {
		d, _ := decimal.NewFromString(redis_db_idx)
		var conn *redis.Client
		conn = redis.NewClient(&redis.Options{
			Addr: redis_addr,
			DB:   int(d.IntPart()),
		})
		e = conn.Ping().Err()
		if e != nil {
			panic(e)
		}
		return conn
	})

	s.logger = di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger)

	s.DBStore = NewDBStore()
	s.NatsStore = NewNatsStore()
	s.RedisStore = NewRedisStore()
	s.PlayerStore = NewPlayerStore()
	s.APIStore = api.NewAPIStore()

	return s
}
