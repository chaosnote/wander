package api

import "context"

// 資料庫表單( agent )
type agent struct {
	ID         string
	Level      uint
	Name       string
	APIKey     string
	Category   string
	ThirdParty string
}

type APIBuilder func(setting agent) Ship

type Ship interface {
	Takeout(ctx context.Context, their_uid string) (money float64, e error)
	Putin(ctx context.Context, their_uid string, money float64) (left_money float64, e error)
}
