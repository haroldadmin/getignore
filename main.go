package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/haroldadmin/getignore/app"
	"github.com/haroldadmin/getignore/flags"
	_ "github.com/haroldadmin/getignore/flags"
)

func main() {
	log.SetFlags(0)

	appContext, cancel := context.WithCancel(context.Background())

	interrupts := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app, err := app.Create(appContext, app.GetIgnoreOptions{
			ShouldUpdate: flags.NoUpdate,
			SearchQuery:  flags.SearchQuery,
			Output:       flags.OutputFile,
		})

		if err != nil {
			log.Print(err)
		}

		app.Start(appContext)
		done <- struct{}{}
	}()

	select {
	case <-interrupts:
		log.Print("Received interrupt")
		cancel()
	case <-done:
	}
}
