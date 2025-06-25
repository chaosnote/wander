package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/chaosnote/wander/data_center/internal"
	"github.com/chaosnote/wander/utils"
)

type Ship_AAAA struct {
	logger *zap.Logger
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
		logger: di.MustGet(internal.LOGGER_API, "AAAA").(*zap.Logger),
	}

	const msg = "NewShipAAAA"
	ship.logger.Debug(msg, zap.Any("setting", setting))

	return ship
}
