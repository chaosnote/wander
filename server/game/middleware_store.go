package game

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosnote/wander/utils"
	"go.uber.org/zap"
)

type MiddlewareStore interface {
	Logging(next http.Handler) http.Handler
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

//-----------------------------------------------

func NewMiddlewareStore() MiddlewareStore {
	var di = utils.GetDI()

	return &middleware_store{
		logger: di.MustGet(LOGGER_SYSTEM).(*zap.Logger),
	}
}
