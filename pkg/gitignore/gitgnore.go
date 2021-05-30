package gitignore

import (
	"errors"
	"path/filepath"

	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/haroldadmin/getignore/pkg/fs"
	"github.com/sahilm/fuzzy"
)

var (
	ErrUnimplemented   = errors.New("unimplemented")
	ErrInvalidWorktree = errors.New("invalid-worktree")
	ErrReadRepoDir     = errors.New("error-read-gitignore-repo-dir")
	ErrNotFound        = errors.New("not-found")
)

var logger = log.WithField("name", "gitignore")

type GitIgnoreFile struct {
	Name string
	Path string
}

type GitIgnoreService interface {
	Get(name string) (GitIgnoreFile, error)
	GetAll() []GitIgnoreFile
	Search(query string) ([]GitIgnoreFile, error)
}

func Create(repository *git.Repository) (GitIgnoreService, error) {
	service := &gitIgnoreService{
		repo: repository,
	}

	err := service.initialize()
	if err != nil {
		return nil, err
	}

	return service, nil
}

type gitIgnoreService struct {
	repo       *git.Repository
	gitIgnores []GitIgnoreFile
}

func (g *gitIgnoreService) initialize() error {
	logger.Infof("initializing GitIgnoreService")

	worktree, err := g.repo.Worktree()
	if err != nil {
		message := "invalid repository worktree"
		logger.Errorf("%s: %v", message, err)
		return ErrInvalidWorktree
	}

	logger.Debugf("reading repo recursively for .gitignore files")
	repoFilesystem := worktree.Filesystem
	files, err := fs.ReadDirRecursively(repoFilesystem)
	if err != nil {
		message := "failed to read dir of gitignore repository"
		logger.Errorf("%s: %v", message, err)
		return ErrReadRepoDir
	}

	gitIgnores := make([]GitIgnoreFile, 0, len(files))
	for _, file := range files {
		ext := filepath.Ext(file.Path)
		if ext != ".gitignore" {
			continue
		}

		gitIgnores = append(gitIgnores, GitIgnoreFile{
			Name: file.Name,
			Path: file.Path,
		})
	}

	logger.Infof("found %d gitignore files", len(gitIgnores))
	for _, f := range gitIgnores {
		logger.Infof("%s (%s)", f.Name, f.Path)
	}

	g.gitIgnores = gitIgnores

	return nil
}

func (g *gitIgnoreService) Get(name string) (GitIgnoreFile, error) {
	logger.Infof("getting file %q", name)

	for _, gitignore := range g.gitIgnores {
		if gitignore.Name == name {
			return gitignore, nil
		}
	}

	logger.Infof("%s not found", name)
	return GitIgnoreFile{}, ErrNotFound
}

func (g *gitIgnoreService) GetAll() []GitIgnoreFile {
	logger.Infof("getting all gitignore files")

	duplicate := make([]GitIgnoreFile, len(g.gitIgnores))
	copy(duplicate, g.gitIgnores)

	return duplicate
}

func (g *gitIgnoreService) Search(query string) ([]GitIgnoreFile, error) {
	logger.Infof("searching gitignore files for %q", query)

	var gitIgnoresDataSource GitIgnores = g.gitIgnores
	matches := fuzzy.FindFrom(query, gitIgnoresDataSource)

	if matches.Len() == 0 {
		logger.Infof("found no matches for %q", query)
		return nil, ErrNotFound
	}

	logger.Infof("found %d matches", matches.Len())
	results := make([]GitIgnoreFile, matches.Len())
	for index, match := range matches {
		results[index] = g.gitIgnores[match.Index]
	}

	return results, nil
}
