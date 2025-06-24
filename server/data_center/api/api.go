package api

import (
	"database/sql"

	"github.com/chaosnote/wander/utils"
)

type APIStore interface {
	APIGet(agent_id string) Ship
}

type api_store struct {
	utils.LogStore

	db *sql.DB

	ship_store map[string]Ship
}

func (s *api_store) init() {
	query := "SELECT * FROM `agent` ;"
	rows, e := s.db.Query(query)
	if e != nil {
		panic(e)
	}
	defer rows.Close()

	// 客制化
	source := map[string]APIBuilder{}
	source["AAAA"] = NewShipAAAA
	// 初始資訊
	for rows.Next() {
		var tmp_agent agent
		e = rows.Scan(
			&tmp_agent.ID,
			&tmp_agent.Level,
			&tmp_agent.Name,
			&tmp_agent.APIKey,
			&tmp_agent.Category,
			&tmp_agent.ThirdParty,
		)
		if e != nil {
			panic(e)
		}
		s.ship_store[tmp_agent.ID] = source[tmp_agent.Name](tmp_agent)
	}
}

func (s *api_store) APIGet(agent_id string) Ship {
	return s.ship_store[agent_id]
}

//-----------------------------------------------

func NewAPIStore() APIStore {
	di := utils.GetDI()
	store := &api_store{
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),
		db:       di.MustGet(utils.SERVICE_MARIADB).(*sql.DB),

		ship_store: map[string]Ship{},
	}
	store.init()

	return store
}
