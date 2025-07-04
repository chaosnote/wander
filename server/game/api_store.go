package game

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/errs"
	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type APIStore interface {
	Login(params member.ReqLogin) (output member.ResLogin, e error)
	Logout(params map[string]any)
}

type api_store struct {
	logger *zap.Logger
}

func (s *api_store) Login(params member.ReqLogin) (output member.ResLogin, e error) {
	const msg = "Login"

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
		s.logger.Error(msg, zap.Error(e))
		e = errs.E10001.Error()
		return
	}

	var body struct {
		api.HttpResponse
		Content member.ResLogin
	}

	e = json.Unmarshal(res.Body(), &body)
	if e != nil {
		s.logger.Error(msg, zap.Error(e))
		e = errs.E00001.Error()
		return
	}

	if body.Code != api.HttpStatusOK {
		e = fmt.Errorf(body.Code)
		return
	}

	output = body.Content

	return
}

func (s *api_store) Logout(params map[string]any) {
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

func NewAPIStore() APIStore {
	var di = utils.GetDI()

	return &api_store{
		logger: di.MustGet(LOGGER_SYSTEM).(*zap.Logger),
	}
}
