package datacenter

import (
	"database/sql"
	"time"

	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type DBStore interface {
	FindUserByID(uid string) (user member.User, e error)
	UpdateUserLastIPByID(uid, client_ip string) (e error)
	UpsertUser(agent_id, their_uname, their_ugrant string, their_uid int64) (uid string, e error)
}

type db_store struct {
	utils.LogStore

	db *sql.DB
}

func (s *db_store) FindUserByID(uid string) (user member.User, e error) {
	query := "SELECT * FROM `user_list` WHERE `ID` = ?"
	row := s.db.QueryRow(query, uid)
	e = row.Scan(
		&user.ID,
		&user.LastIP,
		&user.AgentID,
		&user.TheirUID,
		&user.TheirUName,
		&user.TheirUGrant,
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

func (s *db_store) UpdateUserLastIPByID(uid, client_ip string) (e error) {
	query := "UPDATE `user_list` SET `LastIP` = ? WHERE `ID` = ? "
	_, e = s.db.Exec(
		query,
		client_ip,
		uid,
	)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12003.Error()
		return
	}
	return
}

func (s *db_store) UpsertUser(agent_id, their_uname, their_ugrant string, their_uid int64) (uid string, e error) {
	row := s.db.QueryRow("CALL upsert_user(?, ?, ?, ?, ?) ;", agent_id, their_uid, their_uname, their_ugrant, time.Now().UTC().Format(time.DateTime))
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
		LogStore: di.MustGet(SERVICE_LOGGER).(utils.LogStore),
		db:       di.MustGet(SERVICE_MARIADB).(*sql.DB),
	}
}
