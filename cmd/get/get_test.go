package get_test

import (
	"testing"

	"github.com/haroldadmin/getignore/cmd/get"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("it should have a usage line", func(t *testing.T) {
		usage := get.GetCmd.Use
		assert.NotEmpty(t, usage)
	})
}
