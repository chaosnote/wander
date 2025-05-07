package game

import (
	"github.com/chaosnote/melody"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/member"
)

func (gs *store) handleConnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	session.Set(model.UID, player)
	gs.game_impl.PlayerJoin(player, session)
	gs.addSession(player.UID, session)
}

func (gs *store) handleDisconnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	gs.game_impl.PlayerExit(player, session)
	gs.rmSession(player.UID)
	gs.Logout(map[string]any{model.UID: player.UID})
}

func (gs *store) handleMessageBinary(session *melody.Session, message []byte) {
	content, exists := session.Get(model.UID)
	if !exists {
		return
	}
	player := content.(member.Player)
	gs.game_impl.PlayerMessageBinary(player, session, message)
}

func (gs *store) handleMessage(session *melody.Session, message []byte) {}
