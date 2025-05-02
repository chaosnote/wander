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

	"github.com/chaosnote/melody"
	"github.com/go-resty/resty/v2"

	"github.com/chaosnote/wander/model/api"
	"github.com/chaosnote/wander/utils"
)

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

func dial() {
	var tmp = *token
	if len(tmp) == 0 {
		tmp = getToken()
	}

	u := url.URL{
		Scheme:   "ws",
		Host:     ":8081",
		Path:     "/ws/00/0000",
		RawQuery: fmt.Sprintf("token=%s", tmp),
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

	logger.Debug(utils.LogFields{"tip": "enter"})

	monkey.HandleConnect(func(s *melody.Session) {
		logger.Debug(utils.LogFields{"tip": "connect"})
	})
	monkey.HandleDisconnect(func(s *melody.Session) {
		logger.Debug(utils.LogFields{"tip": "disconnect"})
		os.Exit(0)
	})
	monkey.HandleMessage(func(s *melody.Session, msg []byte) {
		logger.Debug(utils.LogFields{"msg": string(msg)})
	})
	monkey.HandleMessageBinary(func(s *melody.Session, b []byte) {
		logger.Debug(utils.LogFields{"msg": "binary"})
	})

	for i := 0; i < 1; i++ {
		go dial()
	}

	// go func() {
	// 	time.Sleep(2)
	// 	os.Exit(0)
	// }()

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	monkey.Close()

	logger.Debug(utils.LogFields{"tip": "exit"})
	logger.Flush()
}
