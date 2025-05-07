package member

type ReqLogin struct {
	Token  string `json:"token"`
	GameID string `json:"game_id"`
	IP     string `json:"ip"`
}

type ResLogin struct {
	AgentID string  `json:"agent_id"`
	UID     string  `json:"uid"`
	UName   string  `json:"uname"`
	Wallet  float64 `json:"wallet"` // 玩家當前
}
