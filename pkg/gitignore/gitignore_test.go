package gitignore_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/haroldadmin/getignore/pkg/gitignore"
	"github.com/stretchr/testify/assert"
)

func TestGitignoreService(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if given file is not found", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			_, err = service.Get("test.gitignore")
			assert.Error(t, err)
			assert.True(t, errors.Is(err, gitignore.ErrNotFound))
		})

		t.Run("it should not return an error if the given file is found", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			fileName := "Go.gitignore"
			file, err := service.Get(fileName)

			assert.NoError(t, err)
			assert.Equal(t, fileName, file.Name)
		})

		t.Run("it should be case sensitive with file names", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			fileName := "Go.gitignore"
			file, err := service.Get(fileName)

			assert.NoError(t, err)
			assert.Equal(t, fileName, file.Name)

			fileName = strings.ToLower(fileName)
			file, err = service.Get(fileName)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, gitignore.ErrNotFound))
		})
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Run("it should return all discovered gitignore files", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			all := service.GetAll()
			assert.NotEmpty(t, all)
		})

		t.Run("it should return empty list if there are no gitignore files", func(t *testing.T) {
			repo := emptyTestRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			all := service.GetAll()
			assert.Empty(t, all)
		})
	})

	t.Run("Search", func(t *testing.T) {
		t.Run("it should return list of matches for a query", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			matches, err := service.Search("go")
			assert.NoError(t, err)

			hasMatch := false
			for _, match := range matches {
				if match.Name == "Go.gitignore" {
					hasMatch = true
					break
				}
			}

			assert.True(t, hasMatch)
		})

		t.Run("it should return an error if no matches are found", func(t *testing.T) {
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			_, err = service.Search("golang")
			assert.Error(t, err)
			assert.True(t, errors.Is(err, gitignore.ErrNotFound))
		})
	})
}

func testRepository(t *testing.T) *git.Repository {
	t.Helper()

	fs := osfs.New(filepath.Join("testdata", "gitignore"))
	storage := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	repo, err := git.Open(storage, fs)
	if err != nil {
		t.Fatalf("failed to create test repo: %v", err)
	}

	return repo
}

func emptyTestRepository(t *testing.T) *git.Repository {
	t.Helper()

	fs := memfs.New()
	storage := memory.NewStorage()
	repo, err := git.Init(storage, fs)
	if err != nil {
		t.Fatalf("failed to create empty test repo: %v", err)
	}

	return repo
}
