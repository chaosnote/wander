package game

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/chaosnote/melody"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/model/subj"
	"github.com/chaosnote/wander/utils"
)

// 遊戲開發者實做
type GameImpl interface {
	Start(logger utils.LogStore)
	Close()

	PlayerJoin(player member.Player)
	PlayerMessageBinary(player member.Player, pack *message.GameMessage)
	PlayerExit(player member.Player)
}

//-----------------------------------------------

// 底層開放功能
type GameStore interface {
	Start() // 遊戲啟動、此階段會呼叫開發者的對應函式
	Close() // 遊戲關閉、此階段會呼叫開發者的對應函式

	RegisterHandler(provider GameImpl) // 註冊開發者實作函式

	SendGamePack(player member.Player, action string, payload []byte) (e error) // 發送遊戲封包
}

//-----------------------------------------------

type store struct {
	utils.LogStore

	APIStore
	SessionStore
	WalletStore

	game_impl GameImpl
	mel_store *melody.Melody
}

func (s *store) RegisterHandler(provider GameImpl) {
	s.game_impl = provider
}

func (s *store) SendGamePack(player member.Player, action string, payload []byte) (e error) {
	session, ok := s.SessionGet(player.UID)
	if !ok {
		e = errs.E30002.Error()
		s.Info(utils.LogFields{"error": e.Error()})
		return
	}

	var pack = &message.GameMessage{
		Type:      message.GameMessage_RESPONSE,
		Action:    action,
		Payload:   payload,
		Timestamp: utils.UTCUnix(),
	}

	var content []byte
	content, e = proto.Marshal(pack)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00005.Error()
		return
	}
	e = session.WriteBinary(content)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E31001.Error()
		return
	}
	return
}

//-----------------------------------------------

func NewGameStore() GameStore {
	flag.Parse()

	s := &store{}

	var e error
	var di = utils.GetDI()

	di.Set(utils.SERVICE_LOGGER, func() any {
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
		conn.Subscribe(utils.Subject(*GAME_ID, subj.PLAYER_KICK, "*"), s.HandlePlayerKick)
		return conn
	})

	di.Set(utils.SERVICE_REDIS, func() any {
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

	s.LogStore = di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore)

	s.APIStore = NewAPIStore()
	s.SessionStore = NewSessionStore()
	s.WalletStore = NewWalletStore()

	s.mel_store = melody.New()

	return s
}

//-----------------------------------------------

func (s *store) Start() {
	var e error

	s.mel_store.HandleConnect(s.handleConnect)
	s.mel_store.HandleDisconnect(s.handleDisconnect)
	s.mel_store.HandleMessage(s.handleMessage)
	s.mel_store.HandleMessageBinary(s.handleMessageBinary)

	middleware := NewMiddlewareStore()
	// router
	router := mux.NewRouter()
	router.Use(middleware.Logging)

	// [TODO] 第三方驗證功能<Google...>
	// router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {})

	sub := router.PathPrefix("/ws").Subrouter()
	sub.HandleFunc(fmt.Sprintf("/%s/%s", group_id, *GAME_ID), s.HandleGameConn).Queries("token", "{token:[a-zA-Z0-9]{128}}")

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

	go func() {
		e = http.ListenAndServe(game_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()

	var di = utils.GetDI()
	s.game_impl.Start(di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore))
}

func (s *store) Close() {
	s.game_impl.Close()
	s.mel_store.Close()

	var di = utils.GetDI()
	di.MustGet(utils.SERVICE_MARIADB).(*sql.DB).Close()
	di.MustGet(utils.SERVICE_NATS).(*nats.Conn).Close()
	di.MustGet(utils.SERVICE_REDIS).(*redis.Client).Close()
	di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore).Flush()
}
