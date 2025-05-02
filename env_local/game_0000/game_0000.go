package main

import (
	"github.com/chaosnote/melody"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type Game0000 struct {
	utils.LogStore
}

func (g *Game0000) Start(logger utils.LogStore) {
	// 遊戲啟動
	g.LogStore = logger
	g.Debug(utils.LogFields{"tip": "game_start"})
}

func (g *Game0000) Close() {
	// 遊戲關閉
}

func (g *Game0000) PlayerJoin(player member.Player, session *melody.Session) {
	// 玩家上線
	// 玩家基礎資訊
	session.Write([]byte(player.UID))
}

func (g *Game0000) PlayerMessageBinary(player member.Player, session *melody.Session, message []byte) {
	// 玩家封包
}

func (g *Game0000) PlayerExit(player member.Player, session *melody.Session) {
	// 玩家離線
}
