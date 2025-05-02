package api

const HttpStatusOK = "OK"

// HttpResponse 用於回應 http 請求資訊
type HttpResponse struct {
	Code    string // 值解析過程，是否出現錯誤
	Content any    //
}
