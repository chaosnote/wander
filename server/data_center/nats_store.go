package datacenter

import (
	"errors"
	"time"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/subj"
	"github.com/chaosnote/wander/utils"
	"github.com/nats-io/nats.go"
)

type NatsStore interface {
	PlayerKick(player *member.Player, err_id error) (e error)
}

type nats_store struct {
	utils.LogStore

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
		LogStore: di.MustGet(SERVICE_LOGGER).(utils.LogStore),
		conn:     di.MustGet(SERVICE_NATS).(*nats.Conn),
	}
}
