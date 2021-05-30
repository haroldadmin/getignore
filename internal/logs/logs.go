package logs

import (
	"fmt"

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

func CreateLogger(name string) *log.Entry {
	logger := log.WithField("name", name)
	logger.Level = GetLogLevel()

	return logger
}

func setupLogLevel(verbose, veryVerbose bool) {
	logLevel = log.ErrorLevel
	if verbose {
		fmt.Println("Setting info level")
		logLevel = log.InfoLevel
	}

	if veryVerbose {
		fmt.Println("Setting debug level")
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
}

func setupLogHandler() {
	log.SetHandler(cli.Default)
}
