package logs

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

type LogConfig struct {
	Verbose bool
}

func SetupLogs(config LogConfig) {
	setupLogLevel(config.Verbose)
	setupLogHandler()
}

func setupLogLevel(verbose bool) {
	logLevel := log.ErrorLevel
	if verbose {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
	}

func setupLogHandler() {
	log.SetHandler(cli.Default)
}
