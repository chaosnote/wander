package datacenter

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"

	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type DCStore interface {
	Start()
	Close()
}

//-----------------------------------------------

type dc_store struct {
	utils.LogStore

	db_store     *sql.DB
	nats_store   *nats.Conn
	player_store map[string]*member.Player
}

//-----------------------------------------------

func NewDCStore(logger utils.LogStore) DCStore {
	return &dc_store{
		LogStore:     logger,
		player_store: make(map[string]*member.Player),
	}
}

//-----------------------------------------------

func (ds *dc_store) Start() {
	var e error

	// 例 : "user:password@tcp(ip)?parseTime=true/dbname"
	cmd := fmt.Sprintf(`%s:%s@tcp(%s)/%s?parseTime=true`, db_user, db_pw, db_addr, db_name)
	ds.Debug(utils.LogFields{"cmd": cmd})
	ds.db_store, e = sql.Open("mysql", cmd)
	if e != nil {
		panic(e)
	}
	e = ds.db_store.Ping()
	if e != nil {
		panic(e)
	}
	ds.db_store.SetMaxOpenConns(100)          // Limit to N open connections
	ds.db_store.SetMaxIdleConns(10)           // Keep up to N idle connections
	ds.db_store.SetConnMaxLifetime(time.Hour) // Reuse connections for at most N

	// nats
	ds.nats_store, e = nats.Connect(fmt.Sprintf("nats://%s", nats_addr))
	if e != nil {
		panic(e)
	}

	utils.RSAInit("./asset/_rsa/pem.key", 512, true)

	// router
	router := mux.NewRouter()
	router.Use(ds.loggingMiddleware)

	sub := router.PathPrefix("/guest").Subrouter()
	sub.Use(ds.guestMiddleware)
	sub.HandleFunc(`/new`, ds.guestNewHandler).Methods(http.MethodGet)

	sub = router.PathPrefix("/player").Subrouter()
	sub.HandleFunc(`/login`, ds.apiLoginHandler).Methods(http.MethodPost)
	sub.HandleFunc(`/logout`, ds.apiLogoutHandler).Methods(http.MethodPost)

	e = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, e := route.GetPathTemplate()
		if e != nil {
			return e
		}
		ds.Debug(utils.LogFields{"path": template})
		return nil
	})

	if e != nil {
		panic(e)
	}

	ds.Debug(utils.LogFields{"dc_addr": dc_addr})
	go func() {
		e = http.ListenAndServe(dc_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (ds *dc_store) Close() {
	ds.db_store.Close()
	ds.nats_store.Close()
}
