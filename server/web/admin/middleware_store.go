package admin

import (
	"time"

	"github.com/chaosnote/wander/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MiddlewareStore interface {
	Logging() gin.HandlerFunc
}

type middleware_store struct {
}

func (s *middleware_store) Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var di = utils.GetDI()
		logger := di.MustGet(LOGGER_GIN).(*zap.Logger)

		start_time := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		latency := time.Since(start_time)

		param := map[string]interface{}{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"latency":    latency,
		}

		const msg = "Request Record"
		if len(c.Errors) > 0 {
			param["errors"] = c.Errors.Errors()
		}

		logger.Info(msg, zap.Any("param", param))

	}
}

//-----------------------------------------------

func NewMiddlewareStore() MiddlewareStore {
	return &middleware_store{}
}
