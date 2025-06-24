package admin

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/chaosnote/wander/utils"
	"github.com/chaosnote/wander/web/admin/controlller"
	"github.com/gorilla/mux"
)

type AdminStore interface {
	Start()
	Close()
}

type admin_store struct {
	utils.LogStore
}

func (s *admin_store) Start() {
	flag.Parse()

	var di = utils.GetDI()
	var e error

	di.Set(utils.SERVICE_LOGGER, func(args ...interface{}) any {
		file_name := args[0].(string)
		var logger utils.LogStore
		var log_path = filepath.Join("./logs/admin", file_name)
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

	s.LogStore = di.MustGet(utils.SERVICE_LOGGER, "base").(utils.LogStore)

	var middleware = NewMiddlewareStore()

	// router
	router := mux.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.JSON)
	//
	router.HandleFunc("/version", controlller.Version)
	//
	// sub := router.PathPrefix("/version").Subrouter()
	// sub

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

	admin_addr := ":9090"
	s.Debug(utils.LogFields{"admin_addr": admin_addr})
	go func() {
		e = http.ListenAndServe(admin_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (s *admin_store) Close() {}

//-----------------------------------------------

func NewAdminStore() AdminStore {
	return &admin_store{}
}
