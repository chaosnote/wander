package main

import (
	"github.com/chaosnote/wander/game"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/utils"
	"google.golang.org/protobuf/proto"

	_ "github.com/looplab/fsm"

	"idv/chris/protobuf"
	"idv/chris/server/model"
)

type Game0000 struct {
	utils.LogStore
	game.GameStore
}

func (g *Game0000) Start(logger utils.LogStore) {
	// 遊戲啟動
	g.LogStore = logger
	g.Debug(utils.LogFields{"tip": "game_start"})
}

func (g *Game0000) Close() {
	// 遊戲關閉
	g.Debug(utils.LogFields{"tip": "game_close"})
}

func (g *Game0000) PlayerJoin(player member.Player) {
	g.Debug(utils.LogFields{"join": player})

	// 玩家上線<資訊>
	// ∟ 選擇(角色=>對應位置)

	var m model.Game0000Model
	g.RecordLoad(player.UID, m)

	content := &protobuf.Init{
		Player: &protobuf.Player{
			Name:   player.UName,
			Wallet: player.Wallet,
		},
	}

	payload, e := proto.Marshal(content)
	if e != nil {
		g.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00005.Error()
		return
	}
	g.GameStore.SendGamePack(player, "init", payload)
}

func (g *Game0000) PlayerMessageBinary(player member.Player, pack *message.GameMessage) {
	// 玩家封包
}

func (g *Game0000) PlayerExit(player member.Player) {
	// 玩家離線
}
