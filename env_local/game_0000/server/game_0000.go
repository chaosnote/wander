package main

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	_ "github.com/looplab/fsm"

	"github.com/chaosnote/wander/game"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/utils"

	"idv/chris/protobuf"
	"idv/chris/server/model"
)

// 單人遊戲

type Game0000 struct {
	logger *zap.Logger
	game.GameStore
}

func (g *Game0000) Start() {
	var di = utils.GetDI()
	g.logger = di.MustGet(game.LOGGER_GAME, "0000").(*zap.Logger)
	// 遊戲啟動
	g.logger.Debug("game_start")
}

func (g *Game0000) Close() {
	// 遊戲關閉
	g.logger.Debug("game_close")
	g.logger.Sync()
}

func (g *Game0000) PlayerJoin(player member.Player) {
	const msg = "PlayerJoin"
	player_logger := utils.GetDI().MustGet(game.LOGGER_GAME, player.UID).(*zap.Logger)
	defer player_logger.Sync()

	player_logger.Debug(msg, zap.Any("player", player))

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
		player_logger.Error(msg, zap.Error(e))
		e = errs.E00005.Error()
		return
	}
	g.GameStore.SendGamePack(player, protobuf.Action_INIT.String(), payload)
}

func (g *Game0000) PlayerMessageBinary(player member.Player, pack *message.GameMessage) {
	const msg = "PlayerJoin"
	player_logger := utils.GetDI().MustGet(game.LOGGER_GAME, player.UID).(*zap.Logger)
	defer player_logger.Sync()

	// 處理玩家封包
	// ∟ 斷點

	player_logger.Debug(msg, zap.String("action", pack.Action))

	switch pack.Action {
	case protobuf.Action_BET.String():
		g.GameStore.SendGamePack(player, protobuf.Action_BET.String(), nil)
	case protobuf.Action_COMPLETE.String():
		g.GameStore.SendGamePack(player, protobuf.Action_COMPLETE.String(), nil)
	default:
		player_logger.Error(msg, zap.Error(fmt.Errorf("unknow action %s", pack.Action)))
	}
}

func (g *Game0000) PlayerExit(player member.Player) {
	// const msg = "PlayerJoin"
	// player_logger := utils.GetDI().MustGet(game.LOGGER_GAME, player.UID).(*zap.Logger)

	// 玩家離線
	// 是否觸發自動結束
	// ∟ 是 觸發下一個狀態
	// ∟ 否 儲存斷點
}
