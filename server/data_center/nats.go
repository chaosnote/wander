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

func (ds *dc_store) pubPlayerKick(player *member.Player, err_id error) (e error) {
	for i := 0; i < model.NATS_ATTEMPT; i++ {
		_, e = ds.nats_store.Request(utils.Subject(player.GameID, subj.PLAYER_KICK, player.UID), []byte(err_id.Error()), model.NATS_TIMEOUT)
		if e == nil {
			return
		}
		if !errors.Is(e, nats.ErrTimeout) {
			time.Sleep(time.Second)
		}
	}
	return
}
