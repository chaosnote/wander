package model

import "time"

const (
	UID = "UID"

	NATS_TIMEOUT = 3 * time.Second // N 秒無回應，則視為超時
	NATS_ATTEMPT = 5               // 發時錯誤的嘗試次數
)
