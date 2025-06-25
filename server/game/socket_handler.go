package game

import (
	"fmt"

	"github.com/chaosnote/melody"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/model/message"
)

func (s *store) handleConnect(session *melody.Session) {
	player := session.MustGet(model.KEY_UID).(member.Player)

	s.logger.Info(fmt.Sprintf("uid(%s) etner", player.UID))

	session.Set(model.KEY_UID, player)
	s.SessionAdd(player.UID, session)

	s.Takeout(WalletSetting{
		Player:    player,
		Diff:      int(player.Wallet),
		AfterDiff: int(player.Wallet),
	})

	s.game_impl.PlayerJoin(player)
}

func (s *store) handleDisconnect(session *melody.Session) {
	player := session.MustGet(model.KEY_UID).(member.Player)
	s.logger.Info(fmt.Sprintf("uid(%s) leave", player.UID))
	s.game_impl.PlayerExit(player)
	s.SessionRemove(player.UID)

	s.Putin(WalletSetting{
		Player:     player,
		BeforeDiff: int(player.Wallet),
		Diff:       0 - int(player.Wallet),
		AfterDiff:  0,
	})

	s.Logout(map[string]any{
		model.KEY_UID:    player.UID,
		model.KEY_WALLET: player.Wallet,
	})
}

func (s *store) handleMessageBinary(session *melody.Session, source []byte) {
	const msg = "handleMessageBinary"
	content, exists := session.Get(model.KEY_UID)
	if !exists {
		return
	}
	player := content.(member.Player)

	pack := &message.GameMessage{}
	e := proto.Unmarshal(source, pack)
	if e != nil {
		s.logger.Info(msg, zap.Error(e))
		pack.Reset()
		e = errs.E00005.Error()
		*pack.Error = e.Error()
		return
	}

	s.game_impl.PlayerMessageBinary(player, pack)
}

func (s *store) handleMessage(session *melody.Session, message []byte) {}
