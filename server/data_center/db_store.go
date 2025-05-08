package datacenter

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type DBStore interface {
	FindUserByID(agent_id, uid string) (user member.User, e error)
	InsertUser(agent_id, their_uname, their_ugrant string, their_uid, wallet int64) (uid string, e error)
	UpdateUserLastIPByID(agent_id, uid, client_ip string) (e error)
}

type db_store struct {
	utils.LogStore

	db *sql.DB
}

func (s *db_store) FindUserByID(agent_id, uid string) (user member.User, e error) {
	query := fmt.Sprintf("SELECT * FROM `agent_%s_user` WHERE `ID` = ?", agent_id)
	row := s.db.QueryRow(query, uid)
	e = row.Scan(
		&user.ID,
		&user.LastIP,
		&user.TheirUID,
		&user.TheirUName,
		&user.TheirUGrant,
		&user.Wallet,
		&user.CreatedAt,
		&user.ModifiedAt,
	)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12002.Error()
		return
	}
	return
}

func (s *db_store) UpdateUserLastIPByID(agent_id, uid, client_ip string) (e error) {
	query := fmt.Sprintf("UPDATE `agent_%s_user` SET `LastIP` = ?, `ModifiedAt` = ? WHERE `ID` = ? ", agent_id)
	_, e = s.db.Exec(
		query,

		client_ip,
		time.Now().UTC().Format(time.DateTime),
		uid,
	)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12003.Error()
		return
	}
	return
}

func (s *db_store) InsertUser(agent_id, their_uname, their_ugrant string, their_uid, wallet int64) (uid string, e error) {
	query := fmt.Sprintf("INSERT INTO `agent_%s_user` (`TheirUID`, `TheirUName`, `TheirUGrant`, `Wallet`, `CreatedAt`, `ModifiedAt`) VALUES (?, ?, ?, ?, ?, ?);", agent_id)
	create_at := time.Now().UTC().Format(time.DateTime)
	_, e = s.db.Exec(
		query,

		their_uid,
		their_uname,
		their_ugrant,
		wallet,
		create_at,
		create_at,
	)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12001.Error()
		return
	}

	query = fmt.Sprintf("SELECT ID FROM `agent_%s_user` WHERE TheirUID = ? ;", agent_id)
	row := s.db.QueryRow(query, their_uid)
	e = row.Scan(&uid)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12001.Error()
		return
	}
	return
}

//-----------------------------------------------

func NewDBStore() DBStore {
	var di = utils.GetDI()

	return &db_store{
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),
		db:       di.MustGet(utils.SERVICE_MARIADB).(*sql.DB),
	}
}
