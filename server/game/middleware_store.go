package game

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosnote/wander/utils"
)

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
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

//-----------------------------------------------

func NewMiddlewareStore() MiddlewareStore {
	var di = utils.GetDI()

	return &middleware_store{
		LogStore: di.MustGet(SERVICE_LOGGER).(utils.LogStore),
	}
}
