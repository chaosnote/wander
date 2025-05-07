package game

import (
	"strings"

	"github.com/nats-io/nats.go"
)

func (s *store) HandlePlayerKick(msg *nats.Msg) {
	defer msg.Respond(nil)

	var uid = strings.Split(msg.Subject, ".")[2]
	session, ok := s.SessionGet(uid)
	if !ok {
		return
	}
	session.CloseWithMsg(msg.Data)
}
