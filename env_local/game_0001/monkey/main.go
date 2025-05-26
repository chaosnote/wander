package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"google.golang.org/protobuf/proto"

	"github.com/chaosnote/melody"
	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/model/message"
	"github.com/chaosnote/wander/utils"

	"idv/chris/protobuf"
)

// go run . -step 2

type CustomHTTPResponse struct {
	api.HttpResponse
	Content map[string]string
}

var (
	logger = utils.NewConsoleLogger(1)
	monkey = melody.NewMonkey()
)

var (
	token = flag.String("token", "", "token")
	step  = flag.Int("step", 1, "運行多少猴子")
)

func getToken() (token string) {
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
	logger.Debug(utils.LogFields{"token": token})
	return
}

func dial(value string) {
	u := url.URL{
		Scheme:   "ws",
		Host:     ":8081",
		Path:     "/ws/00/0000",
		RawQuery: fmt.Sprintf("token=%s", value),
	}
	logger.Debug(utils.LogFields{"url": u.String()})

	e := monkey.Dial(
		u,
		map[string]any{},
	)
	if e != nil {
		logger.Error(e)
	}
}

func main() {
	flag.Parse()

	var tmp_token = *token
	if len(tmp_token) == 0 {
		tmp_token = getToken()
	}

	logger.Debug(utils.LogFields{"tip": "enter"})

	monkey.HandleConnect(func(s *melody.Session) {
		logger.Debug(utils.LogFields{"tip": "connect"})
	})
	monkey.HandleDisconnect(func(s *melody.Session) {
		logger.Debug(utils.LogFields{"tip": "disconnect"})
	})
	monkey.HandleMessage(func(s *melody.Session, msg []byte) {
		logger.Debug(utils.LogFields{"msg": string(msg)})
	})
	monkey.HandleMessageBinary(func(s *melody.Session, b []byte) {
		pack := &message.GameMessage{}
		e := proto.Unmarshal(b, pack)
		if e != nil {
			panic(e)
		}
		logger.Debug(utils.LogFields{"action": pack.GetAction()})

		content := &protobuf.Init{}
		e = proto.Unmarshal(pack.GetPayload(), content)
		if e != nil {
			panic(e)
		}
		player := content.GetPlayer()
		logger.Debug(utils.LogFields{"player": player})
	})

	for i := 0; i < *step; i++ {
		go dial(tmp_token)
	}

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	monkey.Close()

	logger.Debug(utils.LogFields{"tip": "exit"})
	logger.Flush()
}
