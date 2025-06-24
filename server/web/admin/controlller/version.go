package controlller

import "net/http"

func Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("0.0.0"))
}
