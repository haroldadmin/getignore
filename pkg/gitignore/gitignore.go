package gitignore

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/haroldadmin/getignore/internal/logs"
	"github.com/haroldadmin/getignore/pkg/fs"
	"github.com/sahilm/fuzzy"
)

var (
	ErrUnimplemented   = errors.New("unimplemented")
	ErrInvalidWorktree = errors.New("invalid-worktree")
	ErrReadRepoDir     = errors.New("error-read-gitignore-repo-dir")
	ErrNotFound        = errors.New("not-found")
	ErrInvalidFile     = errors.New("invalid-file")
	ErrReadFile        = errors.New("failed-to-read-file")
	ErrCopyFile        = errors.New("failed-to-copy-file")
	ErrFlushChanges    = errors.New("failed-to-flush-changes")
)

type GitIgnoreFile struct {
	Name string
	Path string
}

type GitIgnoreService interface {
	Get(name string) (GitIgnoreFile, error)
	GetAll() []GitIgnoreFile
	Search(query string) ([]GitIgnoreFile, error)
	Write(file GitIgnoreFile, destFs billy.Filesystem) error
	Append(file GitIgnoreFile, destFs billy.Filesystem) error
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
	logger := logs.CreateLogger("gitignore.init")
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
		logger.Debugf("%s (%s)", f.Name, f.Path)
	}

	g.gitIgnores = gitIgnores
	return nil
}

func (g *gitIgnoreService) Get(name string) (GitIgnoreFile, error) {
	logger := logs.CreateLogger("gitignore.get")
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
	logger := logs.CreateLogger("gitignore.getall")
	logger.Infof("getting all gitignore files")

	duplicate := make([]GitIgnoreFile, len(g.gitIgnores))
	copy(duplicate, g.gitIgnores)

	return duplicate
}

func (g *gitIgnoreService) Search(query string) ([]GitIgnoreFile, error) {
	logger := logs.CreateLogger("gitignore.search")
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

func (g *gitIgnoreService) Write(file GitIgnoreFile, destFs billy.Filesystem) error {
	logger := logs.CreateLogger("gitignore.write")
	worktree, err := g.repo.Worktree()
	if err != nil {
		message := "invalid worktree"
		logger.Infof("%s: %v", message, err)
		return ErrInvalidWorktree
	}

	repoFs := worktree.Filesystem
	destFile, err := destFs.Create(".gitignore")
	if err != nil {
		message := "failed to create/truncate .gitignore file"
		logger.Errorf("%s: %v", message, err)
		return ErrInvalidFile
	}
	defer destFile.Close()

	srcFile, err := repoFs.Open(file.Path)
	if err != nil {
		message := "failed to open source .gitignore file"
		logger.Errorf("%s: %v", message, err)
		return ErrInvalidFile
	}
	defer srcFile.Close()

	destWriter := bufio.NewWriter(destFile)
	destWriter.WriteString("\n\n")
	destWriter.WriteString("# ----- Contents written by getignore -----")
	destWriter.WriteString("\n\n")
	if err := destWriter.Flush(); err != nil {
		message := fmt.Sprintf("failed to flush changes to %q", destFile.Name())
		logger.Errorf("%s: %v", message, err)
		return ErrFlushChanges
	}

	bytesWritten, err := io.Copy(destFile, srcFile)
	if err != nil {
		message := fmt.Sprintf("failed to copy from %q to %q", srcFile.Name(), destFile.Name())
		logger.Errorf("%s: %v", message, err)
		return ErrCopyFile
	}

	logger.Infof("copied %d bytes", bytesWritten)
	return nil
}

func (g *gitIgnoreService) Append(file GitIgnoreFile, destFs billy.Filesystem) error {
	logger := logs.CreateLogger("gitignore.append")
	worktree, err := g.repo.Worktree()
	if err != nil {
		message := "invalid worktree"
		logger.Infof("%s: %v", message, err)
		return ErrInvalidWorktree
	}

	repoFs := worktree.Filesystem
	destFile, err := destFs.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		message := "failed to open .gitignore file"
		logger.Errorf("%s: %v", message, err)
		return ErrInvalidFile
	}
	defer destFile.Close()

	srcFile, err := repoFs.Open(file.Path)
	if err != nil {
		message := "failed to open source .gitignore file"
		logger.Errorf("%s: %v", message, err)
		return ErrInvalidFile
	}
	defer srcFile.Close()

	destWriter := bufio.NewWriter(destFile)
	destWriter.WriteString("\n\n")
	destWriter.WriteString("# ----- Contents written by getignore -----")
	destWriter.WriteString("\n\n")

	srcScanner := bufio.NewScanner(srcFile)
	lineCount := 0
	for srcScanner.Scan() {
		line := fmt.Sprintln(srcScanner.Text())
		destWriter.WriteString(line)
		lineCount++
	}

	if err := destWriter.Flush(); err != nil {
		message := fmt.Sprintf("failed to flush changes to %q", destFile.Name())
		logger.Errorf("%s: %v", message, err)
		return ErrFlushChanges
	}

	logger.Infof("appended %d lines to %q", lineCount, destFile.Name())
	return nil
}
