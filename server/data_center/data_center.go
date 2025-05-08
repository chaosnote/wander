package datacenter

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"

	"github.com/chaosnote/wander/utils"
)

type DCStore interface {
	Start()
	Close()
}

//-----------------------------------------------

type store struct {
	utils.LogStore

	DBStore
	NatsStore
	PlayerStore
}

//-----------------------------------------------

func (s *store) Start() {
	flag.Parse()

	utils.RSAInit("./asset/_rsa/pem.key", 512, true)

	var di = utils.GetDI()
	var e error

	di.Set(utils.SERVICE_LOGGER, func() any {
		var logger utils.LogStore
		var log_path = filepath.Join("./logs/data_center")
		switch *LOG_MODE {
		case 0:
			logger = utils.NewConsoleLogger(1)
		case 2:
			logger = utils.NewMixLogger(log_path, 2)
		default:
			logger = utils.NewFileLogger(log_path, 1)
		}
		return logger
	})

	di.Set(utils.SERVICE_MARIADB, func() any {
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
		db.SetMaxOpenConns(100)          // Limit to N open connections
		db.SetMaxIdleConns(10)           // Keep up to N idle connections
		db.SetConnMaxLifetime(time.Hour) // Reuse connections for at most N

		return db
	})

	di.Set(utils.SERVICE_NATS, func() any {
		var conn *nats.Conn
		conn, e = nats.Connect(fmt.Sprintf("nats://%s", nats_addr))
		if e != nil {
			panic(e)
		}
		return conn
	})

	s.LogStore = di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore)
	s.DBStore = NewDBStore()
	s.NatsStore = NewNatsStore()
	s.PlayerStore = NewPlayerStore()

	middleware := NewMiddlewareStore()
	// router
	router := mux.NewRouter()
	router.Use(middleware.Logging)

	// [TODO]儲值 (未處理)
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
		s.Debug(utils.LogFields{"path": template})
		return nil
	})
	if e != nil {
		panic(e)
	}

	s.Debug(utils.LogFields{"dc_addr": dc_addr})
	go func() {
		e = http.ListenAndServe(dc_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (ds *store) Close() {
	var di = utils.GetDI()
	di.MustGet(utils.SERVICE_MARIADB).(*sql.DB).Close()
	di.MustGet(utils.SERVICE_NATS).(*nats.Conn).Close()
	di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore).Flush()
}

//-----------------------------------------------

func NewDCStore() DCStore {
	return &store{}
}
