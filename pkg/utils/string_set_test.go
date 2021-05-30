package utils_test

import (
	"testing"

	"github.com/haroldadmin/getignore/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringSet(t *testing.T) {
	t.Run("it should return length 0 when empty", func(t *testing.T) {
		set := utils.NewSet()
		assert.True(t, set.Length() == 0)
	})

	t.Run("it should prevent duplicate elements", func(t *testing.T) {
		set := utils.NewSet()
		set.Add("a").Add("b").Add("a")

		assert.True(t, set.Length() == 2)
	})

	t.Run("it should report presence of elements correctly", func(t *testing.T) {
		set := utils.NewSet()
		set.Add("a")

		assert.True(t, set.Contains("a"))
		assert.False(t, set.Contains("b"))
	})

	t.Run("it should remove elements correctly", func(t *testing.T) {
		set := utils.NewSet()
		set.Add("a").Add("b").Remove("b")

		assert.True(t, set.Length() == 1)
		assert.True(t, set.Contains("a"))
		assert.False(t, set.Contains("b"))
	})

	t.Run("it should return keys list with all elements of the set", func(t *testing.T) {
		set := utils.NewSet()
		keys := set.Keys()
		assert.Empty(t, keys)

		set.Add("a").Add("a")
		keys = set.Keys()
		assert.True(t, len(keys) == 1)

		set.Add("b").Add("c").Add("d")
		keys = set.Keys()

		for _, key := range []string{"a", "b", "c", "d"} {
			assert.Contains(t, keys, key)
		}
	})
}
