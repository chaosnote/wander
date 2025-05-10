package api

import (
	"context"

	"github.com/chaosnote/wander/utils"
)

type Ship_AAAA struct {
	utils.LogStore
}

func (s *Ship_AAAA) Takeout(ctx context.Context, their_uid string) (money float64, e error) {
	return
}

func (s *Ship_AAAA) Putin(ctx context.Context, their_uid string, money float64) (left_money float64, e error) {
	left_money = money
	return
}

//-----------------------------------------------

func NewShipAAAA(setting agent) Ship {
	di := utils.GetDI()
	ship := &Ship_AAAA{
		LogStore: di.MustGet(utils.SERVICE_LOGGER).(utils.LogStore),
	}

	ship.Debug(utils.LogFields{"setting": setting})

	return ship
}
