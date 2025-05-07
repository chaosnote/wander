package game

import (
	"github.com/chaosnote/melody"
	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/utils"
)

func (s *store) handleConnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	session.Set(model.UID, player)
	s.Debug(utils.LogFields{"join_uid": player.UID})
	s.game_impl.PlayerJoin(player)
	s.SessionAdd(player.UID, session)
}

func (s *store) handleDisconnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	s.game_impl.PlayerExit(player)
	s.SessionRemove(player.UID)
	s.Logout(map[string]any{model.UID: player.UID})
}

func (s *store) handleMessageBinary(session *melody.Session, source []byte) {
	content, exists := session.Get(model.UID)
	if !exists {
		return
	}
	player := content.(member.Player)

	pack := &message.GameMessage{}
	e := proto.Unmarshal(source, pack)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		pack.Reset()
		e = errs.E00005.Error()
		*pack.Error = e.Error()
		return
	}

	s.game_impl.PlayerMessageBinary(player, pack)
}

func (s *store) handleMessage(session *melody.Session, message []byte) {}
