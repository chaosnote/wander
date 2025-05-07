package game

import (
	"sync"

	"github.com/chaosnote/melody"
)

var session_mu sync.Mutex

type SessionStore interface {
	SessionAdd(uid string, session *melody.Session)
	SessionRemove(uid string)
	SessionGet(uid string) (session *melody.Session, ok bool)
}

type session_store struct {
	pool map[string]*melody.Session
}

func (s *session_store) SessionAdd(uid string, session *melody.Session) {
	session_mu.Lock()
	defer session_mu.Unlock()

	if _, ok := s.pool[uid]; ok {
		return
	}
	s.pool[uid] = session
}

func (s *session_store) SessionGet(uid string) (session *melody.Session, ok bool) {
	session_mu.Lock()
	defer session_mu.Unlock()

	session, ok = s.pool[uid]
	return
}

func (s *session_store) SessionRemove(uid string) {
	session_mu.Lock()
	defer session_mu.Unlock()

	delete(s.pool, uid)
}

//-----------------------------------------------

func NewSessionStore() SessionStore {
	return &session_store{
		pool: map[string]*melody.Session{},
	}
}
