package datacenter

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (s *store) HandleGuestNew(w http.ResponseWriter, r *http.Request) {
	const (
		agent_id     = "57b4866772254df0b157e7966a7c12d2"
		their_ugrant = "1"
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
	uid, e = s.UpsertUser(agent_id, their_uname, their_ugrant, their_uid)
	if e != nil {
		return
	}

	var res = map[string]any{}
	res["token"], e = utils.RSAEncode([]byte(uid))
	if e != nil {
		s.Info(utils.LogFields{"error": e.Error()}) // 思考 DS.Error
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

func (s *store) HandleAPILogin(w http.ResponseWriter, r *http.Request) {
	var e error

	defer func() {
		if e != nil {
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

	player.UID = string(plaintext)
	user, e := s.FindUserByID(player.UID) // 是否為已註冊的使用者
	if e != nil {
		return
	}
	player.UName = user.TheirUName

	old_player, ok := s.PlayerAdd(player) // old_player {true:玩家增加成功:false:玩家增加失敗}
	if !ok {
		e = errs.E13001.Error()
		s.PlayerKick(old_player, e)
		return
	}

	e = s.UpdateUserLastIPByID(player.UID, player.IP)
	if e != nil {
		return
	}

	output, e := json.Marshal(api.HttpResponse{Code: api.HttpStatusOK, Content: player.ResLogin})
	if e != nil {
		s.Error(e)
		e = errs.E00001.Error()
		return
	}
	w.Write(output)
}

func (s *store) HandleAPILogout(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	utils.HttpRequestJSONUnmarshal(r.Body, &body)
	s.PlayerRemove(body[model.UID])
}
