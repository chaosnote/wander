package datacenter

import (
	"sync"

	"github.com/chaosnote/wander/model/member"
)

var player_mu sync.Mutex

func (ds *dc_store) addPlayer(new_player *member.Player) (old_player *member.Player, add_suc bool) {
	player_mu.Lock()
	defer player_mu.Unlock()

	var ok bool
	old_player, ok = ds.player_store[new_player.UID]
	if ok {
		add_suc = false
		return
	}
	add_suc = true
	ds.player_store[new_player.UID] = new_player
	return
}

func (ds *dc_store) rmPlayer(uid string) {
	player_mu.Lock()
	defer player_mu.Unlock()

	delete(ds.player_store, uid)
}

func (ds *dc_store) getPlayer(uid string) (player *member.Player, ok bool) {
	player_mu.Lock()
	defer player_mu.Unlock()

	player, ok = ds.player_store[uid]

	return
}
