package game

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (s *store) HandleGameConn(w http.ResponseWriter, r *http.Request) {
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
			s.Info(utils.LogFields{"error": e.Error()})

			conn, inrupt := s.mel_store.Upgrader.Upgrade(w, r, w.Header())
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

	player.ResLogin, e = s.Login(player.ReqLogin)
	if e != nil {
		return
	}

	var keys = make(map[string]any)
	keys[model.KEY_UID] = player

	e = s.mel_store.HandleRequestWithKeys(w, r, keys)
	if e != nil {
		s.Error(e)
		e = errs.E30006.Error()
		return
	}

}
