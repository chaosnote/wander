package member

import "time"

type User struct {
	ID          uint      `json:"id"`
	LastIP      string    `json:"last_ip"`
	TheirUID    string    `json:"their_uid"`
	TheirUName  string    `json:"their_uname"`
	TheirUGrant string    `json:"their_ugrant"`
	Wallet      float64   `json:"wallet"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}
