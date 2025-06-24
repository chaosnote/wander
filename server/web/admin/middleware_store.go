package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosnote/wander/utils"
)

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
	JSON(next http.Handler) http.Handler
}

type middleware_store struct {
}

func (s *middleware_store) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := utils.GetDI().MustGet(utils.SERVICE_LOGGER, "router").(utils.LogStore)

		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()

		duration := endTime.Sub(startTime)

		logger.Info(utils.LogFields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"duration":  fmt.Sprintf("%v", duration),
			"client_ip": utils.ParseIP(r),
		})
	})
}

func (s *middleware_store) JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

//-----------------------------------------------

func NewMiddlewareStore() MiddlewareStore {
	return &middleware_store{}
}
