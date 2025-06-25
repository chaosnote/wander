package datacenter

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/subj"
	"github.com/chaosnote/wander/utils"
)

type NatsStore interface {
	PlayerKick(player *member.Player, err_id error) (e error)
}

type nats_store struct {
	logger *zap.Logger

	conn *nats.Conn
}

func (s *nats_store) PlayerKick(player *member.Player, err_id error) (e error) {
	for i := 0; i < model.NATS_ATTEMPT; i++ {
		_, e = s.conn.Request(utils.Subject(player.GameID, subj.PLAYER_KICK, player.UID), []byte(err_id.Error()), model.NATS_TIMEOUT)
		if e == nil {
			return
		}
		if !errors.Is(e, nats.ErrTimeout) {
			time.Sleep(time.Second)
		}
	}
	return
}

//-----------------------------------------------

func NewNatsStore() NatsStore {
	var di = utils.GetDI()

	return &nats_store{
		logger: di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger),
		conn:   di.MustGet(internal.SERVICE_NATS).(*nats.Conn),
	}
}
