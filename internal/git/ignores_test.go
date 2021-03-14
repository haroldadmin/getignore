package git_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/haroldadmin/getignore/internal/git"
)

func TestIgnores(t *testing.T) {
	t.Run("should find .gitignore files only one-level deep", func(t *testing.T) {
		repoDir := createRepoDir(t)
		createGitignore(t, "root.gitignore", repoDir)
		createGitignore(t, "deep.gitignore", filepath.Join(repoDir, "deepDir"))

		ignores := git.Ignores{RepoDir: repoDir}
		if err := ignores.FindIgnores(); err != nil {
			t.Errorf("expected no errors when finding ignores, got: %v", err)
		}

		if numIgnores := ignores.Len(); numIgnores != 1 {
			t.Errorf("expected one gitignore file only, got: %v", numIgnores)
		}

		if exists := checkExists(t, ignores.GitIgnores, "root.gitignore"); !exists {
			t.Error("expected discovered gitignores to contain root.gitignore, but it wasn't there")
		}

		if exists := checkExists(t, ignores.GitIgnores, "deep.gitignore"); exists {
			t.Error("expected discovered gitignores to not contain deep.gitignore, but it was there")
		}
	})

	t.Run("should limit search results to the number supplied", func(t *testing.T) {
		repoDir := createRepoDir(t)
		gitignoreFiles := []string{"a.gitignore", "b.gitignore", "c.gitignore", "d.gitignore", "e.gitignore"}
		for _, file := range gitignoreFiles {
			createGitignore(t, file, repoDir)
		}

		ignores := git.Ignores{RepoDir: repoDir}
		if err := ignores.FindIgnores(); err != nil {
			t.Errorf("expected no errors when finding ignores, got: %v", err)
		}

		results := ignores.SearchIgnores("gitignore", 1)
		if numResults := len(results); numResults != 1 {
			t.Errorf("expected one result, got %v", numResults)
		}
	})

	t.Run("should return .gitignore files in search results if they are present and match the query", func(t *testing.T) {
		repoDir := createRepoDir(t)
		gitignoreFiles := []string{"a.gitignore", "b.gitignore", "c.gitignore", "d.gitignore", "e.gitignore"}
		for _, file := range gitignoreFiles {
			createGitignore(t, file, repoDir)
		}

		ignores := git.Ignores{RepoDir: repoDir}
		if err := ignores.FindIgnores(); err != nil {
			t.Errorf("expected no errors when finding ignores, got: %v", err)
		}

		results := ignores.SearchIgnores("a.gitignore", 5)
		if numResults := len(results); numResults != 1 {
			t.Errorf("expected one result, got %v", numResults)
		}

		if exists := checkExists(t, results, "a.gitignore"); !exists {
			t.Error("expected results to contain a.gitignore, but it wasn't there")
		}
	})
}

func createRepoDir(t *testing.T) string {
	t.Helper()

	repoDir := filepath.Join(t.TempDir(), "temp-repo")
	err := os.MkdirAll(repoDir, os.ModePerm)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create repo dir: %v", err))
	}

	t.Cleanup(func() {
		os.RemoveAll(repoDir)
	})

	return repoDir
}

func createGitignore(t *testing.T, name string, dir string) string {
	t.Helper()

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatal(fmt.Errorf("parent directory %v invalid: %v", dir, err))
	}

	gitIgnorePath := filepath.Join(dir, name)
	f, err := os.Create(gitIgnorePath)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create %v: %v", gitIgnorePath, err))
	}

	defer f.Close()

	t.Cleanup(func() {
		os.Remove(gitIgnorePath)
	})

	return f.Name()
}

func checkExists(t *testing.T, ignores []git.GitIgnoreFile, name string) bool {
	t.Helper()

	for _, i := range ignores {
		if i.Name == name {
			return true
		}
	}

	return false
}
