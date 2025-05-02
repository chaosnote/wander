package game

import (
	"strings"

	"github.com/nats-io/nats.go"
)

func (gs *game_store) handlePlayerKick(msg *nats.Msg) {
	defer msg.Respond(nil)

	var uid = strings.Split(msg.Subject, ".")[2]
	session, ok := gs.getSession(uid)
	if !ok {
		return
	}
	session.CloseWithMsg(msg.Data)
}
