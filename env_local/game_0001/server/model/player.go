package model

import (
	"github.com/chaosnote/wander/model/member"
)

type GamePlayer struct {
	Player member.Player

	Mode mode
}

func NewGamePlayer(player member.Player) *GamePlayer {
	return &GamePlayer{
		Player: player,

		Mode: MODE_0,
	}
}
