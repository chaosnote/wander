package member

import "time"

type User struct {
	ID          uint      `json:"id"`
	LastIP      string    `json:"last_ip"`
	AgentID     string    `json:"agent_id"`
	TheirUID    string    `json:"their_uid"`
	TheirUName  string    `json:"their_uname"`
	TheirUGrant string    `json:"their_ugrant"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}
