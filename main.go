package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/haroldadmin/getignore/app"
	"github.com/haroldadmin/getignore/flags"
	_ "github.com/haroldadmin/getignore/flags"
)

func main() {
	initLogger()
	interrupts := initInterrupts()

	appContext, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}, 1)

	go func() {
		app, err := app.Create(appContext, app.GetIgnoreOptions{
			ShouldUpdate: flags.NoUpdate,
			SearchQuery:  flags.SearchQuery,
			Output:       flags.OutputFile,
		})

		if err != nil {
			log.WithError(err).Fatal("Failed to initialize getignore")
		}

		app.Start(appContext)
		done <- struct{}{}
	}()

	select {
	case <-interrupts:
		log.Debug("Received interrupt")
		cancel()
	case <-done:
	}
}

func initLogger() {
	logLevel := log.WarnLevel
	if flags.Verbose {
		logLevel = log.DebugLevel
	}

	log.SetHandler(cli.Default)
	log.SetLevel(logLevel)
}

func initInterrupts() chan os.Signal {
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupts
}
