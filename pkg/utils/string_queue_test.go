package utils_test

import (
	"testing"

	"github.com/haroldadmin/getignore/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringQueue(t *testing.T) {
	t.Run("it should report length 0 when queue is empty", func(t *testing.T) {
		q := utils.NewQueue()
		assert.Equal(t, q.Length(), 0)
	})

	t.Run("it should add an element to the queue using Add()", func(t *testing.T) {
		q := utils.NewQueue()
		q.Add("first")

		assert.Equal(t, q.Length(), 1)
	})

	t.Run("it should return the first element of the queue using RemoveFirst()", func(t *testing.T) {
		q := utils.NewQueue()
		q.Add("first").Add("second").Add("third")

		first, err := q.RemoveFirst()

		assert.NoError(t, err)
		assert.Equal(t, "first", first)
	})

	t.Run("it should return an error with RemoveFirts() when queue is empty", func(t *testing.T) {
		q := utils.NewQueue()

		first, err := q.RemoveFirst()

		assert.Error(t, err)
		assert.Empty(t, first)
	})
}
