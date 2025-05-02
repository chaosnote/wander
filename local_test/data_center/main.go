package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	dc "github.com/chaosnote/wander/data_center"
	"github.com/chaosnote/wander/utils"
)

func main() {
	logger := utils.NewConsoleLogger(1)
	server := dc.NewDCStore(logger)

	logger.Debug(utils.LogFields{"tip": "server starting"})

	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q

	logger.Debug(utils.LogFields{"tip": "server closing"})
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Close()

	logger.Debug(utils.LogFields{"tip": "server stop"})
	logger.Flush()
}
