package logs_test

import (
	"testing"

	"github.com/apex/log"
	"github.com/haroldadmin/getignore/internal/logs"
	"github.com/stretchr/testify/assert"
)

func TestSetupLogs(t *testing.T) {
	t.Run("it should set log level to info for verbose config", func(t *testing.T) {
		logs.SetupLogs(logs.LogConfig{
			Verbose: true,
		})

		assert.Equal(t, logs.GetLogLevel(), log.InfoLevel)
	})

	t.Run("it should set log level to debug for very verbose config", func(t *testing.T) {
		logs.SetupLogs(logs.LogConfig{
			VeryVerbose: true,
		})

		assert.Equal(t, logs.GetLogLevel(), log.DebugLevel)
	})
}
