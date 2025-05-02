package model

import "time"

const (
	UID = "UID"

	PACKET_TIMEOUT = 3 * time.Second // N 秒無回應，則視為超時
	PACKET_ATTEMPT = 5               // 發時錯誤的嘗試次數
)
