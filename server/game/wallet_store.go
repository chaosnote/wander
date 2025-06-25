package game

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type WalletAction uint

const (
	TAKEOUT  WalletAction = iota // 取出
	DEDUCT                       // 支出
	ROLLBACK                     // 支出回滾
	INCOME                       // 收入
	PUTIN                        // 歸還
	REFRESH                      // 更新
)

type WalletSetting struct {
	member.Player

	RoundID    sql.NullString
	BeforeDiff int
	Diff       int
	AfterDiff  int
}

type Wallet struct {
	WalletSetting

	ActionType          WalletAction
	TransactionDatetime time.Time
}

type WalletStore interface {
	Takeout(wallet_setting WalletSetting) (e error)
	Putin(wallet_setting WalletSetting) (e error)
}

type wallet_store struct {
	logger *zap.Logger
	db     *sql.DB
}

func (s *wallet_store) store(wallet Wallet) (e error) {

	query := fmt.Sprintf(
		"INSERT INTO `agent_%s_wallet` (`ID`, `TheirUID`, `GameID`, `RoundID`, `BeforeDiff`, `Diff`, `AfterDiff`, `ActionType`, `TransactionDatetime`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		wallet.AgentID,
	)

	our_uid, _ := decimal.NewFromString(wallet.UID)

	_, e = s.db.Exec(
		query,

		int(our_uid.IntPart()),
		wallet.TheirUID,
		wallet.GameID,
		wallet.RoundID,
		wallet.BeforeDiff,
		wallet.Diff,
		wallet.AfterDiff,
		wallet.ActionType,
		wallet.TransactionDatetime,
	)
	if e != nil {
		return
	}
	return
}

func (s *wallet_store) Takeout(wallet_setting WalletSetting) (e error) {
	const msg = "Takeout"
	wallet := Wallet{
		WalletSetting:       wallet_setting,
		ActionType:          TAKEOUT,
		TransactionDatetime: time.Now().UTC(),
	}
	e = s.store(wallet)
	if e != nil {
		s.logger.Error(msg, zap.Error(e))
		e = errs.E32001.Error()
	}
	return
}

func (s *wallet_store) Deduct(wallet_setting WalletSetting) {}

func (s *wallet_store) Rollback(wallet_setting WalletSetting) {}

func (s *wallet_store) Income(wallet_setting WalletSetting) {}

func (s *wallet_store) Putin(wallet_setting WalletSetting) (e error) {
	const msg = "Putin"

	wallet := Wallet{
		WalletSetting:       wallet_setting,
		ActionType:          PUTIN,
		TransactionDatetime: time.Now().UTC(),
	}
	e = s.store(wallet)
	if e != nil {
		s.logger.Error(msg, zap.Error(e))
		e = errs.E32002.Error()
	}
	return
}

func (s *wallet_store) Refresh(agent_id string, wallet_setting WalletSetting) {}

//-----------------------------------------------

func NewWalletStore() WalletStore {
	var di = utils.GetDI()

	return &wallet_store{
		logger: di.MustGet(LOGGER_SYSTEM).(*zap.Logger),
		db:     di.MustGet(SERVICE_MARIADB).(*sql.DB),
	}
}
