package search_test

import (
	"testing"

	"github.com/haroldadmin/getignore/cmd/search"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("it should have a usage line", func(t *testing.T) {
		usage := search.SearchCmd.Use
		assert.NotEmpty(t, usage)
	})
}
