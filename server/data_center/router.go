package datacenter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"

	"github.com/chaosnote/wander/model"
	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (ds *dc_store) initRouter() {
	router := mux.NewRouter()
	router.Use(ds.loggingMiddleware)

	var e error

	sub := router.PathPrefix("/guest").Subrouter()
	sub.Use(ds.guestMiddleware)
	sub.HandleFunc(`/new`, ds.guestNewHandler).Methods(http.MethodGet)

	sub = router.PathPrefix("/player").Subrouter()
	sub.HandleFunc(`/login`, ds.apiLoginHandler).Methods(http.MethodPost)
	sub.HandleFunc(`/logout`, ds.apiLogoutHandler).Methods(http.MethodPost)

	e = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, e := route.GetPathTemplate()
		if e != nil {
			return e
		}
		ds.Debug(utils.LogFields{"path": template})
		return nil
	})

	if e != nil {
		panic(e)
	}

	ds.Debug(utils.LogFields{"dc_addr": dc_addr})
	go func() {
		e = http.ListenAndServe(dc_addr, router)
		if e != nil && e != http.ErrServerClosed {
			panic(e)
		}
	}()
}

func (ds *dc_store) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()

		duration := endTime.Sub(startTime)

		ds.Info(utils.LogFields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"duration":  fmt.Sprintf("%v", duration),
			"client_ip": utils.ParseIP(r),
		})
	})
}

func (ds *dc_store) guestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if use_guest != "1" {

			ds.Info(utils.LogFields{
				"method":    r.Method,
				"use_guest": use_guest,
				"client_ip": utils.ParseIP(r),
			})

			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func (ds *dc_store) guestNewHandler(w http.ResponseWriter, r *http.Request) {
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
	row := ds.db_store.QueryRow("CALL upsert_user(?, ?, ?, ?, ?) ;", agent_id, their_uid, their_uname, their_ugrant, time.Now().UTC().Format(time.DateTime))
	e = row.Scan(&uid)
	if e != nil {
		ds.Info(utils.LogFields{"error": e.Error()}) // [TODO]DS.Error
		e = errs.E12001.Error()
		return
	}

	var res = map[string]any{}
	res["token"], e = utils.RSAEncode([]byte(uid))
	if e != nil {
		ds.Info(utils.LogFields{"error": e.Error()}) // 思考 DS.Error
		e = errs.E12001.Error()
		return
	}

	output, e := json.Marshal(api.HttpResponse{Code: api.HttpStatusOK, Content: res})
	if e != nil {
		ds.Error(e)
		e = errs.E00001.Error()
		return
	}
	w.Write(output)
}

func (ds *dc_store) apiLoginHandler(w http.ResponseWriter, r *http.Request) {
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
		ds.Info(utils.LogFields{"error": e.Error()})
		e = errs.E10004.Error()
		return
	}

	ds.Debug(utils.LogFields{"params": player.ReqLogin})

	var data []byte
	data, e = utils.HexDecode(player.Token)
	if e != nil {
		ds.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00002.Error()
		return
	}

	var plaintext []byte
	plaintext, e = utils.RSADecode(string(data))
	if e != nil {
		ds.Info(utils.LogFields{"error": e.Error()})
		e = errs.E00003.Error()
		return
	}

	player.UID = string(plaintext)
	user, e := ds.findUserByID(player.UID) // 是否為已註冊的使用者
	if e != nil {
		return
	}
	player.UName = user.TheirUName

	_, ok := ds.addPlayer(player) // old_player {true:玩家增加成功:false:玩家增加失敗}
	if !ok {
		// ds.nats_store.Publish(utils.Subject(old_player.GateID, subj.PLAYER_KICK, old_player.UID), nil)
		e = errs.E13001.Error()
		return
	}

	e = ds.updateUserLastIPByID(player.UID, player.IP)
	if e != nil {
		return
	}

	output, e := json.Marshal(api.HttpResponse{Code: api.HttpStatusOK, Content: player.ResLogin})
	if e != nil {
		ds.Error(e)
		e = errs.E00001.Error()
		return
	}
	w.Write(output)
}

func (ds *dc_store) apiLogoutHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	utils.HttpRequestJSONUnmarshal(r.Body, &body)
	ds.rmPlayer(body[model.UID])
}
