package game

import (
	"github.com/chaosnote/melody"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (s *store) handleConnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	session.Set(model.UID, player)
	s.Debug(utils.LogFields{"join_uid": player.UID})
	s.game_impl.PlayerJoin(player, session)
	s.SessionAdd(player.UID, session)
}

func (s *store) handleDisconnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	s.game_impl.PlayerExit(player, session)
	s.SessionRemove(player.UID)
	s.Logout(map[string]any{model.UID: player.UID})
}

func (s *store) handleMessageBinary(session *melody.Session, message []byte) {
	content, exists := session.Get(model.UID)
	if !exists {
		return
	}
	player := content.(member.Player)
	s.game_impl.PlayerMessageBinary(player, session, message)
}

func (s *store) handleMessage(session *melody.Session, message []byte) {}
