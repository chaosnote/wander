package datacenter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosnote/wander/utils"
)

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
	Guest(next http.Handler) http.Handler
}

type middleware_store struct {
	utils.LogStore
}

func (s *middleware_store) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()

		duration := endTime.Sub(startTime)

		s.Info(utils.LogFields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"duration":  fmt.Sprintf("%v", duration),
			"client_ip": utils.ParseIP(r),
		})
	})
}

func (s *middleware_store) Guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if use_guest != "1" {

			s.Info(utils.LogFields{
				"method":    r.Method,
				"use_guest": use_guest,
				"client_ip": utils.ParseIP(r),
			})

			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

//-----------------------------------------------

func NewMiddlewareStore() MiddlewareStore {
	var di = utils.GetDI()

	return &middleware_store{
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),
	}
}
