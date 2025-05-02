package game

import (
	"sync"

	"github.com/chaosnote/melody"
)

var session_mu sync.Mutex

func (gs *game_store) addSession(uid string, session *melody.Session) {
	session_mu.Lock()
	defer session_mu.Unlock()
	// [TODO]
	if _, ok := gs.session_store[uid]; ok {
		return
	}
	gs.session_store[uid] = session
}

func (gs *game_store) getSession(uid string) (session *melody.Session, ok bool) {
	session_mu.Lock()
	defer session_mu.Unlock()

	session, ok = gs.session_store[uid]
	return
}

func (gs *game_store) rmSession(uid string) {
	session_mu.Lock()
	defer session_mu.Unlock()

	delete(gs.session_store, uid)
}
