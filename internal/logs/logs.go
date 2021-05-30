package logs

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var logLevel log.Level

type LogConfig struct {
	Verbose     bool
	VeryVerbose bool
}

func SetupLogs(config LogConfig) {
	setupLogLevel(config.Verbose, config.VeryVerbose)
	setupLogHandler()
}

func GetLogLevel() log.Level {
	return logLevel
}

func setupLogLevel(verbose, veryVerbose bool) {
	logLevel = log.ErrorLevel
	if verbose {
		logLevel = log.InfoLevel
	}

	if veryVerbose {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

func setupLogHandler() {
	log.SetHandler(cli.Default)
}
