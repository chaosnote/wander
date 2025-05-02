package member

type ReqLogin struct {
	Token  string `json:"token"`
	GameID string `json:"game_id"`
	IP     string `json:"ip"`
}

type ResLogin struct {
	UID   string `json:"uid"`
	UName string `json:"uname"`
}
