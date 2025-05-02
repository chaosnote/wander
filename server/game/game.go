package game

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/chaosnote/melody"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/shopspring/decimal"

	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/subj"
	"github.com/chaosnote/wander/utils"
)

// 遊戲開發者實做
type GameImpl interface {
	Start(utils.LogStore)
	Close()

	PlayerJoin(player member.Player, session *melody.Session)
	PlayerMessageBinary(player member.Player, session *melody.Session, message []byte)
	PlayerExit(player member.Player, session *melody.Session)
}

//-----------------------------------------------

type GameStore interface {
	Start()
	Close()
}

//-----------------------------------------------

type game_store struct {
	utils.LogStore

	db_store      *sql.DB
	nats_store    *nats.Conn
	mel_store     *melody.Melody
	redis_store   *redis.Client
	session_store map[string]*melody.Session

	game_impl GameImpl
}

//-----------------------------------------------

func NewGameStore(provider GameImpl) GameStore {
	return &game_store{
		game_impl:     provider,
		session_store: map[string]*melody.Session{},
	}
}

//-----------------------------------------------

func (gs *game_store) Start() {
	flag.Parse()

	var logger utils.LogStore
	var log_path = filepath.Join(log_dir, fmt.Sprintf("game_%s", *GAME_ID))
	switch *LOG_MODE {
	case 0:
		logger = utils.NewConsoleLogger(1)
	case 2:
		logger = utils.NewMixLogger(log_path, 2)
	default:
		logger = utils.NewFileLogger(log_path, 1)
	}
	gs.LogStore = logger

	var e error

	// mariadb
	// 例 : "user:password@tcp(ip)?parseTime=true/dbname"
	cmd := fmt.Sprintf(`%s:%s@tcp(%s)/%s?parseTime=true`, db_user, db_pw, db_addr, db_name)
	gs.Debug(utils.LogFields{"cmd": cmd})
	gs.db_store, e = sql.Open("mysql", cmd)
	if e != nil {
		panic(e)
	}
	e = gs.db_store.Ping()
	if e != nil {
		panic(e)
	}
	// redis
	d, _ := decimal.NewFromString(redis_db_idx)
	gs.redis_store = redis.NewClient(&redis.Options{
		Addr: redis_addr,
		DB:   int(d.IntPart()),
	})
	e = gs.redis_store.Ping().Err()
	if e != nil {
		panic(e)
	}
	// melody
	gs.mel_store = melody.New()
	gs.mel_store.HandleConnect(gs.handleConnect)
	gs.mel_store.HandleDisconnect(gs.handleDisconnect)
	gs.mel_store.HandleMessage(gs.handleMessage)
	gs.mel_store.HandleMessageBinary(gs.handleMessageBinary)
	// nats
	gs.nats_store, e = nats.Connect(fmt.Sprintf("nats://%s", nats_addr))
	if e != nil {
		panic(e)
	}
	gs.nats_store.Subscribe(utils.Subject(*GAME_ID, subj.PLAYER_KICK, "*"), gs.handlePlayerKick)
	// router
	router := mux.NewRouter()
	router.Use(gs.loggingMiddleware)

	// [TODO] 第三方驗證功能<Google...>
	// router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {})

	sub := router.PathPrefix("/ws").Subrouter()
	sub.HandleFunc(fmt.Sprintf("/%s/%s", group_id, *GAME_ID), gs.gameConnHandler).Queries("token", "{token:[a-zA-Z0-9]{128}}")

	e = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, e := route.GetPathTemplate()
		if e != nil {
			return e
		}
		gs.Debug(utils.LogFields{"path": template})
		return nil
	})

	if e != nil {
		panic(e)
	}

	go func() {
		e = http.ListenAndServe(game_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()

	gs.game_impl.Start(logger)
}

func (gs *game_store) Close() {
	gs.game_impl.Close()
	gs.mel_store.Close()
	gs.redis_store.Close()
	gs.db_store.Close()

	gs.Flush()
}
