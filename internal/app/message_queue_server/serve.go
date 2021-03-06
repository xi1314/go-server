// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue_server

import (
	"context"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve() error {
	var (
		consumers []*nsq.Consumer
	)

	redis.Connect()
	database.Connect()

	go func() {
		if ctx, err := RunMessageQueueConsumer(); err != nil {
			log.Fatal(err)
		} else {
			consumers = ctx
		}
	}()

	log.Println("Listening message queue")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	config.Common.Exiting = true

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if len(consumers) > 0 {
		for _, c := range consumers {
			c.Stop()
			_ = c.DisconnectFromNSQD(message_queue.Address)
		}
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("Timeout of 5 seconds.")

	redis.Dispose()
	database.Dispose()

	log.Println("Server exiting")

	return nil
}
