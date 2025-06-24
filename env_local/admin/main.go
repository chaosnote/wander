package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaosnote/wander/web/admin"
)

func main() {
	server := admin.NewAdminStore()
	server.Start()

	log.Println("server starting")

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	log.Println("server closing")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("server stop")
	server.Close()
}
