package model

import "time"

const (
	KEY_AGENT_ID = "agent_id"
	KEY_UID      = "uid"
	KEY_WALLET   = "wallet"

	NATS_TIMEOUT = 3 * time.Second // N 秒無回應，則視為超時
	NATS_ATTEMPT = 5               // 發時錯誤的嘗試次數
)
