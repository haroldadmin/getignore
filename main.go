package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/haroldadmin/getignore/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupts
		log.Debug("received interrupt")
		cancel()
	}()

	cmd.Execute(ctx)
}
