package gitignore_test

import (
	"errors"
	"io/ioutil"
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

	t.Run("Write", func(t *testing.T) {
		t.Run("it should create a .gitignore file if it doesn't exist", func(t *testing.T) {
			repo := testRepository(t)
			destFs := memfs.New()

			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			gitIgnoreFile := service.GetAll()[0]
			err = service.Write(gitIgnoreFile, destFs)
			assert.NoError(t, err)

			files, err := destFs.ReadDir("/")
			assert.NoError(t, err)
			assert.Len(t, files, 1)

			file := files[0]
			assert.Equal(t, file.Name(), ".gitignore")
		})

		t.Run("it should truncate an existing .gitignore file if it exists", func(t *testing.T) {
			destFs := memfs.New()
			f, err := destFs.Create(".gitignore")
			assert.NoError(t, err)
			defer f.Close()

			f.Write([]byte("test-data"))

			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			gitIgnoreFile := service.GetAll()[0]
			err = service.Write(gitIgnoreFile, destFs)
			assert.NoError(t, err)

			contents, err := ioutil.ReadAll(f)
			assert.NoError(t, err)
			assert.NotContains(t, string(contents), "test-data")
		})
	})

	t.Run("Append", func(t *testing.T) {
		t.Run("it should return an error if there is no existing gitignore file", func(t *testing.T) {
			destFs := memfs.New()
			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			gitIgnoreFile := service.GetAll()[0]

			err = service.Append(gitIgnoreFile, destFs)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, gitignore.ErrInvalidFile))
		})

		t.Run("it should append to an existing gitignore file", func(t *testing.T) {
			destFs := memfs.New()
			f, err := destFs.Create(".gitignore")
			assert.NoError(t, err)
			defer f.Close()

			f.Write([]byte("test-data"))
			c, _ := ioutil.ReadAll(f)
			t.Log(string(c))

			repo := testRepository(t)
			service, err := gitignore.Create(repo)
			assert.NoError(t, err)

			gitIgnoreFile := service.GetAll()[0]
			err = service.Append(gitIgnoreFile, destFs)
			assert.NoError(t, err)

			f.Seek(0, 0)
			contents, err := ioutil.ReadAll(f)
			assert.NoError(t, err)

			stringContents := string(contents)
			assert.NoError(t, err)
			assert.True(t, strings.HasPrefix(stringContents, "test-data"))
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
