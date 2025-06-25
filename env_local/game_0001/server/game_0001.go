package main

import (
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

// 多人連線

type Game0001 struct {
	logger *zap.Logger
	game.GameStore
}

func (g *Game0001) Start() {
	var di = utils.GetDI()
	g.logger = di.MustGet(game.LOGGER_GAME, "0001").(*zap.Logger)

	// 遊戲啟動 - 讀取設定檔
	//
	// 建立 loger
	// 建立 (桌/室)
	//
	g.logger.Debug("game_start")
}

func (g *Game0001) Close() {
	// 遊戲關閉
	// ∟ 遊戲室(桌/室)是否回合結束
	// ∟
	g.logger.Debug("game_close")
	g.logger.Sync()
}

func (g *Game0001) PlayerJoin(player member.Player) {
	const msg = "PlayerJoin"
	player_logger := utils.GetDI().MustGet(game.LOGGER_GAME, player.UID).(*zap.Logger)
	defer player_logger.Sync()
	player_logger.Debug(msg, zap.Any("player", player))

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
		player_logger.Error(msg, zap.Error(e))
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
