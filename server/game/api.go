package game

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

func (gs *game_store) login(params member.ReqLogin) (output member.ResLogin, e error) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)

	u := url.URL{
		Scheme:   "http",
		Host:     dc_addr,
		Path:     "/player/login",
		RawQuery: fmt.Sprintf("t=%s", utils.UTCUnixString()),
	}

	res, e := client.R().SetBody(params).Post(u.String())
	if e != nil {
		gs.Error(e)
		e = errs.E10001.Error()
		return
	}

	var body struct {
		api.HttpResponse
		Content member.ResLogin
	}

	e = json.Unmarshal(res.Body(), &body)
	if e != nil {
		gs.Error(e)
		e = errs.E00001.Error()
		return
	}

	if body.Code != api.HttpStatusOK {
		e = fmt.Errorf(body.Code)
		return
	}

	output = body.Content

	gs.Debug(utils.LogFields{"user": output})

	return
}

func (gs *game_store) logout(params map[string]any) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(5)

	u := url.URL{
		Scheme:   "http",
		Host:     dc_addr,
		Path:     "/player/logout",
		RawQuery: fmt.Sprintf("t=%s", utils.UTCUnixString()),
	}

	client.R().SetBody(params).Post(u.String())
}
