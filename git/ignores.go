package git

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/sahilm/fuzzy"
)

type GitIgnoreFile struct {
	Name string
	Path string
}

type Ignores struct {
	RepoDir    string
	gitIgnores []GitIgnoreFile
}

func (i *Ignores) String(index int) string {
	return i.gitIgnores[index].Name
}

func (i *Ignores) Len() int {
	return len(i.gitIgnores)
}

func (i *Ignores) FindIgnores() error {
	pattern := i.RepoDir + string(filepath.Separator) + "*.gitignore"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("Failed to find .gitignore files: %v", err)
	}

	gitIgnores := make([]GitIgnoreFile, 0)
	for _, matchPath := range matches {
		gitIgnore := GitIgnoreFile{
			Name: filepath.Base(matchPath),
			Path: matchPath,
		}
		gitIgnores = append(gitIgnores, gitIgnore)
	}

	log.Printf("Found %d .gitignore files", len(gitIgnores))
	i.gitIgnores = gitIgnores

	return nil
}

func (i *Ignores) SearchIgnores(query string, maxResults int) []GitIgnoreFile {
	matches := fuzzy.FindFrom(query, i)
	results := make([]GitIgnoreFile, 0, matches.Len())

	for matchNumber := 0; matchNumber < maxResults && matchNumber < len(matches); matchNumber++ {
		match := matches[matchNumber]
		results = append(results, i.gitIgnores[match.Index])
	}

	return results
}
