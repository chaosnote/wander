package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaosnote/wander/game"
)

func main() {
	base := game.NewGameStore()
	base.RegisterHandler(&Game0001{GameStore: base})
	base.Start()

	log.Println("server starting")

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	log.Println("server closing")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	base.Close()

	log.Println("server stop")
}
