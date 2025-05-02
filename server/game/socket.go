package game

import (
	"github.com/chaosnote/melody"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/member"
)

func (gs *game_store) handleConnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	session.Set(model.UID, player)
	gs.game_impl.PlayerJoin(player, session)
	gs.addSession(player.UID, session)
}

func (gs *game_store) handleDisconnect(session *melody.Session) {
	player := session.MustGet(model.UID).(member.Player)
	gs.game_impl.PlayerExit(player, session)
	gs.rmSession(player.UID)
	gs.logout(map[string]any{model.UID: player.UID})
}

func (gs *game_store) handleMessageBinary(session *melody.Session, message []byte) {
	content, exists := session.Get(model.UID)
	if !exists {
		return
	}
	player := content.(member.Player)
	gs.game_impl.PlayerMessageBinary(player, session, message)
}

func (gs *game_store) handleMessage(session *melody.Session, message []byte) {}
