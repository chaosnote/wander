package datacenter

import (
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (ds *dc_store) findUserByID(uid string) (user member.User, e error) {
	query := "SELECT * FROM `user_list` WHERE `ID` = ?"
	row := ds.db_store.QueryRow(query, uid)
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
		ds.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12002.Error()
		return
	}
	return
}

func (ds *dc_store) updateUserLastIPByID(uid, client_ip string) (e error) {
	query := "UPDATE `user_list` SET `LastIP` = ? WHERE `ID` = ? "
	_, e = ds.db_store.Exec(
		query,
		client_ip,
		uid,
	)
	if e != nil {
		ds.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12003.Error()
		return
	}
	return
}
