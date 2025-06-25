package game

import (
	"context"
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/model/subj"
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

// 遊戲開發者實做
type GameImpl interface {
	Start()
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

	OutputMongoStore
}

//-----------------------------------------------

type store struct {
	logger *zap.Logger

	APIStore
	SessionStore
	WalletStore
	MongoStore

	game_impl GameImpl
	mel_store *melody.Melody
}

func (s *store) RegisterHandler(provider GameImpl) {
	s.game_impl = provider
}

func (s *store) SendGamePack(player member.Player, action string, payload []byte) (e error) {
	const msg = "SendGamePack"
	session, ok := s.SessionGet(player.UID)
	if !ok {
		e = errs.E30002.Error()
		s.logger.Error(msg, zap.Error(e))
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
		s.logger.Error(msg, zap.Error(e))
		e = errs.E00005.Error()
		return
	}
	e = session.WriteBinary(content)
	if e != nil {
		s.logger.Error(msg, zap.Error(e))
		e = errs.E31001.Error()
		return
	}
	return
}

//-----------------------------------------------

func (s *store) Start() {
	const msg = "Start"

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
		s.logger.Debug(msg, zap.String("path", template))
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

	s.game_impl.Start()
}

func (s *store) Close() {
	s.game_impl.Close()
	s.mel_store.Close()

	var di = utils.GetDI()
	di.MustGet(SERVICE_MARIADB).(*sql.DB).Close()
	di.MustGet(SERVICE_NATS).(*nats.Conn).Close()
	di.MustGet(SERVICE_REDIS).(*redis.Client).Close()
	di.MustGet(SERVICE_MONGO).(*mongo.Client).Disconnect(context.TODO())
	di.MustGet(LOGGER_SYSTEM).(*zap.Logger).Sync()
}

//-----------------------------------------------

func NewGameStore() GameStore {
	flag.Parse()

	s := &store{}

	var e error
	var di = utils.GetDI()

	log_mode := *LOG_MODE
	di.SetShare(LOGGER_SYSTEM, func(...interface{}) any {
		var log_path = filepath.Join(log_dir, fmt.Sprintf("game_%s", *GAME_ID))
		return gen_logger(log_mode, log_path)
	})
	di.Set(LOGGER_GAME, func(args ...interface{}) any {
		var uid = args[0].(string)
		var log_path = filepath.Join(log_dir, *GAME_ID, uid)
		return gen_logger(log_mode, log_path)
	})

	di.SetShare(SERVICE_MARIADB, func(...interface{}) any {
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

	di.SetShare(SERVICE_NATS, func(...interface{}) any {
		var conn *nats.Conn
		conn, e = nats.Connect(fmt.Sprintf("nats://%s", nats_addr))
		if e != nil {
			panic(e)
		}
		conn.Subscribe(utils.Subject(*GAME_ID, subj.PLAYER_KICK, "*"), s.HandlePlayerKick)
		return conn
	})

	di.SetShare(SERVICE_REDIS, func(...interface{}) any {
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

	di.SetShare(SERVICE_MONGO, func(...interface{}) any {
		client_options := options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/", mongo_user, mongo_pw, mongo_addr)).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(5 * time.Minute)

		client, err := mongo.Connect(context.TODO(), client_options)
		if err != nil {
			panic(e)
		}
		return client
	})

	s.logger = di.MustGet(LOGGER_SYSTEM).(*zap.Logger)

	s.APIStore = NewAPIStore()
	s.SessionStore = NewSessionStore()
	s.WalletStore = NewWalletStore()
	s.MongoStore = NewMongoStore()

	s.mel_store = melody.New()

	return s
}
