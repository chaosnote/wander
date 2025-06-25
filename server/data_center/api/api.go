package api

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/utils"
)

type APIStore interface {
	APIGet(agent_id string) Ship
}

type api_store struct {
	logger *zap.Logger

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
		logger: di.MustGet(internal.LOGGER_SYSTEM).(*zap.Logger),
		db:     di.MustGet(internal.SERVICE_MARIADB).(*sql.DB),

		ship_store: map[string]Ship{},
	}
	store.init()

	return store
}
