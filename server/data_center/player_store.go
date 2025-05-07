package datacenter

import (
	"sync"

	"github.com/chaosnote/wander/model/member"
)

var player_mu sync.Mutex

type PlayerStore interface {
	PlayerAdd(new_player *member.Player) (old_player *member.Player, add_suc bool)
	PlayerRemove(uid string)
	PlayerGet(uid string) (player *member.Player, ok bool)
}

type player_store struct {
	pool map[string]*member.Player
}

func (p *player_store) PlayerAdd(new_player *member.Player) (old_player *member.Player, add_suc bool) {
	player_mu.Lock()
	defer player_mu.Unlock()

	var ok bool
	old_player, ok = p.pool[new_player.UID]
	if ok {
		add_suc = false
		return
	}
	add_suc = true
	p.pool[new_player.UID] = new_player
	return
}

func (p *player_store) PlayerRemove(uid string) {
	player_mu.Lock()
	defer player_mu.Unlock()

	delete(p.pool, uid)
}

func (p *player_store) PlayerGet(uid string) (player *member.Player, ok bool) {
	player_mu.Lock()
	defer player_mu.Unlock()

	player, ok = p.pool[uid]

	return
}

//-----------------------------------------------

func NewPlayerStore() PlayerStore {
	return &player_store{
		pool: map[string]*member.Player{},
	}
}
