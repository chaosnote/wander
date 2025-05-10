package datacenter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type uid_queue struct {
	mu    sync.Mutex
	count int
}

var (
	uid_mu    sync.Mutex
	uid_store = make(map[string]*uid_queue)
)

func (s *store) HandleGuestNew(w http.ResponseWriter, r *http.Request) {
	const (
		agent_id     = "92aa258f95834a8bb35f74d5c21787d8"
		their_ugrant = "1"
		wallet       = 100000
	)
	var their_uname = utils.GenSerial("demo_")
	tmp_uid, _ := decimal.NewFromString(strings.Split(utils.GenSerial("uid_"), "_")[1])
	their_uid := tmp_uid.IntPart()

	var e error

	defer func() {
		if e != nil {
			output, _ := json.Marshal(api.HttpResponse{Code: e.Error()})
			w.Write(output)
		}
	}()

	var uid string
	uid, e = s.InsertUser(agent_id, their_uname, their_ugrant, their_uid, wallet)
	if e != nil {
		return
	}

	var res = map[string]any{}
	res["token"], e = utils.RSAEncode([]byte(fmt.Sprintf("%s|%s", agent_id, uid)))
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E12001.Error()
		return
	}

	output, e := json.Marshal(api.HttpResponse{Code: api.HttpStatusOK, Content: res})
	if e != nil {
		s.Error(e)
		e = errs.E00001.Error()
		return
	}
	w.Write(output)
}

func (s *store) HandlePlayerLogin(w http.ResponseWriter, r *http.Request) {
	var e error
	var st = time.Now()

	defer func() {
		if e != nil && !errors.Is(e, errs.E10005.Error()) {
			output, _ := json.Marshal(api.HttpResponse{Code: e.Error()})
			w.Write(output)
		}
	}()

	player := &member.Player{}
	e = utils.HttpRequestJSONUnmarshal(r.Body, &player.ReqLogin)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E10004.Error()
		return
	}
	s.Debug(utils.LogFields{"params": player.ReqLogin})

	var data []byte
	data, e = utils.HexDecode(player.Token)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00002.Error()
		return
	}

	var plaintext []byte
	plaintext, e = utils.RSADecode(string(data))
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00003.Error()
		return
	}

	list := strings.Split(string(plaintext), "|")
	player.AgentID = list[0]
	player.UID = list[1]

	s.Info(utils.LogFields{
		"agent_id": player.AgentID,
		"uid":      player.UID,
	})

	var timeout = 5 * time.Second
	ctx := r.Context()
	ctx_first, cancel_first := context.WithTimeout(ctx, timeout)
	defer cancel_first()

	uid_mu.Lock()
	if _, ok := uid_store[player.UID]; ok {
		uid_store[player.UID].count++
	} else {
		uid_store[player.UID] = &uid_queue{count: 1}
	}
	queue := uid_store[player.UID]
	uid_mu.Unlock()

	queue.mu.Lock()
	defer func() {
		queue.mu.Unlock()
		uid_mu.Lock()
		uid_store[player.UID].count--
		if uid_store[player.UID].count <= 0 {
			delete(uid_store, player.UID)
		}
		uid_mu.Unlock()
	}()

	// 過程中請求端已超時(5 秒)、此時後端仍在運行、則後端資訊與狀態不符
	// 是否為已註冊的使用者
	user, e := s.FindUserByID(player.AgentID, player.UID)
	if e != nil {
		return
	}
	player.UName = user.TheirUName
	player.Wallet = user.Wallet
	player.TheirUID = user.TheirUID

	if !s.BlackNotExisted(player.UID) {
		e = errs.E14001.Error()
		s.Info(utils.LogFields{"error": e.Error(), "uid": player.UID})
		return
	}

	// 已註冊的使用者，但已有連線
	// old_player {true:玩家增加成功:false:玩家增加失敗}
	old_player, ok := s.PlayerAdd(player)
	if !ok {
		e = errs.E13001.Error()
		s.BlackAdd(player.UID, player.Token)
		s.PlayerKick(old_player, e)
		return
	}

	e = s.UpdateUserIPAndWallet(player.AgentID, player.UID, player.IP, 0)
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error(), "uid": player.UID})
		return
	}
	defer func() {
		if e != nil {
			e = s.UpdateUserIPAndWallet(player.AgentID, player.UID, player.IP, player.Wallet)
		}
		if e != nil {
			s.Info(utils.LogFields{"error": e.Error(), "uid": player.UID})
			return
		}
	}()

	var money float64
	money, e = s.APIGet(player.AgentID).Takeout(ctx, player.TheirUID)
	if e != nil {
		if e != nil {
			s.Error(e)
			e = errs.E10001.Error()
			return
		}
	}
	player.Wallet = player.Wallet + money // 本地+額外值
	defer func() {
		if e != nil {
			player.Wallet, e = s.APIGet(player.AgentID).Putin(ctx, player.TheirUID, player.Wallet)
		}
	}()

	output, e := json.Marshal(api.HttpResponse{Code: api.HttpStatusOK, Content: player.ResLogin})
	if e != nil {
		s.Error(e)
		e = errs.E00001.Error()
		return
	}

	select {
	case <-ctx.Done():
		e = errs.E10005.Error()
		s.Info(utils.LogFields{"error": e.Error(), "uid": player.UID})
		return
	case <-ctx_first.Done():
		if time.Since(st) >= timeout {
			e = errs.E10007.Error()
		} else {
			e = errs.E10005.Error()
		}
		s.Info(utils.LogFields{"error": e.Error(), "uid": player.UID})
		return
	default:
	}

	_, e = w.Write(output)
	if e != nil {
		s.Error(e)
		e = errs.E10006.Error()
	}
}

func (s *store) HandlePlayerLogout(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	e := utils.HttpRequestJSONUnmarshal(r.Body, &body)
	if e != nil {
		s.Error(e)
		return
	}

	uid := body[model.KEY_UID].(string)
	player, ok := s.PlayerRemove(uid)
	if !ok {
		s.Info(utils.LogFields{"error": fmt.Sprintf("lost uid(%s)", uid)})
		return
	}

	player.Wallet = body[model.KEY_WALLET].(float64)
	player.Wallet, e = s.APIGet(player.AgentID).Putin(r.Context(), player.TheirUID, player.Wallet)
	if e != nil {
		s.Error(e)
		return
	}
	e = s.UpdateUserIPAndWallet(player.AgentID, player.UID, player.IP, player.Wallet)
	if e != nil {
		s.Error(e)
		return
	}
}
