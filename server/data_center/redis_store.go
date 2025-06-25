package datacenter

import (
	"fmt"
	"time"

	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/utils"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var ErrNil = fmt.Errorf("redis: nil")

type RedisStore interface {
	BlackAdd(uid string, token string)
	BlackNotExisted(uid string) (ok bool)
}

type redis_store struct {
	logger *zap.Logger

	conn *redis.Client
}

func (s *redis_store) BlackAdd(uid string, token string) {
	s.conn.Set(uid, token, 2*time.Second)
}

func (s *redis_store) BlackNotExisted(uid string) (ok bool) {
	const msg = "BlackNotExisted"

	cmd := s.conn.Get(uid)
	e := cmd.Err()
	if e != nil && e.Error() != ErrNil.Error() {
		s.logger.Error(msg, zap.Error(cmd.Err()))
		return
	}
	if len(cmd.Val()) > 0 {
		return
	}
	ok = true
	return
}

//-----------------------------------------------

func NewRedisStore() RedisStore {
	var di = utils.GetDI()

	return &redis_store{
		logger: di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger),
		conn:   di.MustGet(internal.SERVICE_REDIS).(*redis.Client),
	}
}
