package datacenter

import (
	"fmt"
	"time"

	"github.com/chaosnote/wander/utils"
	"github.com/go-redis/redis"
)

var ErrNil = fmt.Errorf("redis: nil")

type RedisStore interface {
	BlackAdd(uid string, token string)
	BlackNotExisted(uid string) (ok bool)
}

type redis_store struct {
	utils.LogStore

	conn *redis.Client
}

func (s *redis_store) BlackAdd(uid string, token string) {
	s.conn.Set(uid, token, 2*time.Second)
}

func (s *redis_store) BlackNotExisted(uid string) (ok bool) {
	cmd := s.conn.Get(uid)
	e := cmd.Err()
	if e != nil && e.Error() != ErrNil.Error() {
		s.Debug(utils.LogFields{"error": cmd.Err()})
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
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),
		conn:     di.MustGet(utils.SERVICE_REDIS).(*redis.Client),
	}
}
