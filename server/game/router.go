package game

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (gs *game_store) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()

		duration := endTime.Sub(startTime)

		gs.Info(utils.LogFields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"duration":  fmt.Sprintf("%v", duration),
			"client_ip": utils.ParseIP(r),
		})
	})
}

func (gs *game_store) gameConnHandler(w http.ResponseWriter, r *http.Request) {
	var e error
	player := member.Player{
		ReqLogin: member.ReqLogin{
			Token:  r.URL.Query().Get("token"),
			GameID: *GAME_ID,
			IP:     utils.ParseIP(r),
		},
	}

	defer func() {
		if e != nil {
			gs.Info(utils.LogFields{"error": e.Error()})

			conn, inrupt := gs.mel_store.Upgrader.Upgrade(w, r, w.Header())
			if inrupt != nil {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			closeMessage := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, e.Error())
			conn.WriteMessage(websocket.CloseMessage, closeMessage)
			conn.Close()
			return
		}
	}()

	player.ResLogin, e = gs.login(player.ReqLogin)
	if e != nil {
		return
	}

	var keys = make(map[string]any)
	keys[model.UID] = player

	e = gs.mel_store.HandleRequestWithKeys(w, r, keys)
	if e != nil {
		gs.Error(e)
		e = errs.E20005.Error()
		return
	}

}
