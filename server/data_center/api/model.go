package api

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
	Takeout(their_uid string) (money float64, e error)
	Putin(their_uid string, money float64) (left_money float64, e error)
}
