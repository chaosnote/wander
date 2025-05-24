package monkey

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/chaosnote/melody"
	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/utils"
)

var (
	game_id      = flag.String("game_id", "0000", "game_id")
	monkey_count = flag.Int("monkey_count", 1, "monkey_count")
	token        = flag.String("token", "", "token")
)

//-----------------------------------------------

type CustomHTTPResponse struct {
	api.HttpResponse
	Content map[string]string
}

//-----------------------------------------------

type MonkeyImpl interface {
}

//-----------------------------------------------

type MonkeyStore interface {
	GetToken() (token string)
	Dial(token string) (e error)
}

type monkey_store struct {
	utils.LogStore
}

func (s *monkey_store) GetToken() (token string) {
	client := resty.New()
	client.SetTimeout(5 * time.Second)

	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   "/guest/new",
	}

	res, e := client.R().Get(u.String())
	if e != nil {
		panic(e)
	}

	output := CustomHTTPResponse{}
	e = json.Unmarshal(res.Body(), &output)
	if e != nil {
		panic(e)
	}
	if output.Code != api.HttpStatusOK {
		panic(output.Code)
	}
	token = output.Content["token"]
	s.Debug(utils.LogFields{"token": token})
	return
}

func (s *monkey_store) Dial(token string) (e error) {
	u := url.URL{
		Scheme:   "ws",
		Host:     ":8081",
		Path:     "/ws/00/0001",
		RawQuery: fmt.Sprintf("token=%s", token),
	}
	s.Debug(utils.LogFields{"url": u.String()})

	e = melody.NewMonkey().Dial(
		u,
		map[string]any{},
	)
	if e != nil {
		s.Error(e)
	}
	return
}

//-----------------------------------------------

func NewMonkeyStore() MonkeyStore {
	return &monkey_store{
		LogStore: utils.NewConsoleLogger(1),
	}
}
