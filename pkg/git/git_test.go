package git

import (
	"context"
	"errors"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	t.Run("it should return an error if context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		storage := memory.NewStorage()
		fs := memfs.New()

		cancel()
		_, err := clone(ctx, storage, fs)
		assert.Error(t, err, "Expected an error")
	})
}

func TestUpdate(t *testing.T) {
	t.Run("it should return an error if repository is invalid", func(t *testing.T) {
		ctx := context.Background()
		uninitializedRepo := git.Repository{}

		err := update(ctx, &uninitializedRepo)
		assert.Error(t, err, "Expected an error")
		assert.True(
			t,
			errors.Is(err, git.ErrIsBareRepository),
			"Expected error to be git.ErrIsBareRepository",
		)
	})

	t.Run("it should return an error if context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		storage := memory.NewStorage()
		fs := memfs.New()

		repo, err := git.Init(storage, fs)
		assert.NoError(t, err, "Expected repository to be initialized successfully")

		cancel()
		err = update(ctx, repo)
		assert.Error(t, err, "Expected an error")
	})
}
