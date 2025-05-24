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

// 多人連線

type Game0001 struct {
	utils.LogStore
	game.GameStore
}

func (g *Game0001) Start(logger utils.LogStore) {
	// 遊戲啟動 - 讀取設定檔
	//
	// 建立 loger
	// 建立 (桌/室)
	//
	g.LogStore = logger
	g.Debug(utils.LogFields{"tip": "game_start"})
}

func (g *Game0001) Close() {
	// 遊戲關閉
	// ∟ 遊戲室(桌/室)是否回合結束
	// ∟
	g.Debug(utils.LogFields{"tip": "game_close"})
}

func (g *Game0001) PlayerJoin(player member.Player) {
	g.Debug(utils.LogFields{"join": player})

	// 玩家狀態
	// ∟ 選擇
	// ∟ 進入選擇

	var m model.GameModel
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

func (g *Game0001) PlayerMessageBinary(player member.Player, pack *message.GameMessage) {
	// 玩家封包
}

func (g *Game0001) PlayerExit(player member.Player) {
	// 玩家離線
}
