package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func ParseIP(r *http.Request) (client_ip string) {
	client_ip = r.Header.Get("X-Forwarded-For")
	if client_ip == "" {
		client_ip = r.Header.Get("X-Real-IP")
	}

	if client_ip == "" {
		client_ip = r.RemoteAddr
		if colon := strings.LastIndex(client_ip, ":"); colon != -1 {
			client_ip = client_ip[:colon]
		}
	} else {
		ips := strings.Split(client_ip, ",")
		client_ip = strings.TrimSpace(ips[0])
	}

	return
}

func HttpRequestJSONUnmarshal(reader io.ReadCloser, output any) (e error) {
	var body []byte
	body, e = io.ReadAll(reader)
	if e != nil {
		return
	}
	defer reader.Close()

	e = json.Unmarshal(body, output)
	if e != nil {
		return
	}
	return
}
