package errs

import "fmt"

type error_id int

func (e error_id) Error() error {
	return fmt.Errorf("E%05d", e)
}

func (e error_id) ErrorWithMessage(msg string) error {
	return fmt.Errorf("E%05d : %s", e, msg)
}

const (
	E00000 error_id = iota //
	E00001                 // 編碼( json )    錯誤
	E00002                 // 編碼( hex )     錯誤
	E00003                 // 編碼( rsa )     錯誤
	E00004                 // 映射( reflect ) 錯誤
	E00005                 // 編碼( proto )   錯誤
)

// DataCenter 資訊中心錯誤
//
// 10000 Http       錯誤起點
// 11000 NatsIO     錯誤起點
// 12000 DB         錯誤起點
// 13000 Model      錯誤起點
// 14000 Redis      錯誤起點

const (
	E10000 error_id = iota + 10000 //
	E10001                         // Http Request 過程出錯
	E10002                         // WEB 服務回應( 非 OK )
	E10003                         // Http Request <參數/格式>錯誤
	E10004                         // Http POST Body 轉換失敗
	E10005                         // Http Response 過程出錯
)

const (
	E12000 error_id = iota + 12000 //
	E12001                         // Insert `user` 過程出錯
	E12002                         // Select `user` 過程出錯
	E12003                         // UPDATE `user` 過程出錯
)

const (
	E13000 error_id = iota + 13000 //
	E13001                         // 已有玩家記錄(重複連線)
)

const (
	E14000 error_id = iota + 14000 //
	E14001                         // 玩家短期內重複連線
)

//
// Game 錯誤
//
// 30000 Http       錯誤起點
// 31000 Socket     錯誤起點
// 32000 DB         錯誤起點

const (
	E30000 error_id = iota + 30000 //
	E30001                         // 已存在相同玩家
	E30002                         // 無指定玩家
	E30003                         // 伺服器關閉前觸發
	E30004                         // http 請求階段出錯
	E30005                         // 玩家閒置過久
	E30006                         // 玩家連線失敗
)

const (
	E31000 error_id = iota + 31000 //
	E31001                         // 封包傳遞失敗
)

const (
	E32000 error_id = iota + 32000 //
	E32001                         // 取值記錄失敗
	E32002                         // 回存記錄失敗
)
