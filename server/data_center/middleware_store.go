package datacenter

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/utils"
)

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
	Guest(next http.Handler) http.Handler
}

type middleware_store struct {
	logger *zap.Logger
}

func (s *middleware_store) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start_time)

		s.logger.Info("Request Record", zap.Any("param", utils.CustomField{
			"method":    r.Method,
			"path":      r.RequestURI,
			"duration":  fmt.Sprintf("%v", latency),
			"client_ip": utils.ParseIP(r),
		}))
	})
}

func (s *middleware_store) Guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if use_guest != "1" {

			s.logger.Info("Request Record", zap.Any("param", utils.CustomField{
				"method":    r.Method,
				"use_guest": use_guest,
				"client_ip": utils.ParseIP(r),
			}))

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
		logger: di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger),
	}
}
