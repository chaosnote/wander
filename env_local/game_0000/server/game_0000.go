package main

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/wander/game"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/utils"

	_ "github.com/looplab/fsm"

	"idv/chris/protobuf"
	"idv/chris/server/model"
)

// 單人遊戲

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

	// 是否有斷點資訊
	// ∟ 是 值反序列為 RoomModel
	// ∟ 否 產生 RoomModel
	//
	// RoomModel 記錄至 GameModel
	//
	var m model.RoomModel
	g.RecordLoad(player.UID, &m)

	// 更新:
	// ∟ 玩家錢包 + 當前[畫面]贏分

	// 測試用[尚未處理]
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
	g.GameStore.SendGamePack(player, protobuf.ActionType_INIT.String(), payload)
}

func (g *Game0000) PlayerMessageBinary(player member.Player, pack *message.GameMessage) {
	// 處理玩家封包
	// ∟ 斷點

	g.Debug(utils.LogFields{"action": pack.Action})
	switch pack.Action {
	case protobuf.ActionType_BET.String():
		g.GameStore.SendGamePack(player, protobuf.ActionType_BET.String(), nil)
	case protobuf.ActionType_COMPLETE.String():
		g.GameStore.SendGamePack(player, protobuf.ActionType_COMPLETE.String(), nil)
	default:
		g.Error(fmt.Errorf("unknow action %s", pack.Action))
	}
}

func (g *Game0000) PlayerExit(player member.Player) {
	// 玩家離線
	// 是否觸發自動結束
	// ∟ 是 觸發下一個狀態
	// ∟ 否 儲存斷點
}
