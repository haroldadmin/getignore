package git

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/haroldadmin/getignore/internal/logs"
)

var (
	ErrInvalidPath    = errors.New("invalid-path")
	ErrGitServiceInit = errors.New("git-service-init-failed")
	ErrChroot         = errors.New("failed-to-chroot")
)

const GitIgnoreRepository string = "https://github.com/github/gitignore"

// CreateOptions contains config parameters for creating a Git repository
type CreateOptions struct {
	RepositoryDir    string
	UpdateRepository bool
}

// Create creates a Git Repository and returns a reference to it
func Create(ctx context.Context, options CreateOptions) (*git.Repository, error) {
	logger := logs.CreateLogger("git.create")
	logger.Infof("creating GitService: %s", options.RepositoryDir)

	repoPath, err := filepath.Abs(options.RepositoryDir)
	if err != nil {
		logger.Errorf("failed to parse absolute path: %v", err)
		return nil, ErrInvalidPath
	}
	repositoryFs := osfs.New(repoPath)

	dotGitFs, err := repositoryFs.Chroot(".git")
	if err != nil {
		logger.Errorf("failed to chroot into .git dir: %v", err)
		return nil, ErrChroot
	}
	dotGitStorage := filesystem.NewStorage(dotGitFs, cache.NewObjectLRUDefault())

	repository, err := initialize(ctx, dotGitStorage, repositoryFs, options)
	if err != nil {
		return nil, ErrGitServiceInit
	}

	return repository, nil
}

// initialize clones or updates the GitIgnore repository
// [storage] must be derived from a filesystem rooted at the `.git` directory
// of the Gitignore repository
// [filesystem] must be rooted at the repository directory.
func initialize(
	ctx context.Context,
	storage storage.Storer,
	filesystem billy.Filesystem,
	options CreateOptions,
) (*git.Repository, error) {
	logger := logs.CreateLogger("git.init")
	logger.Info("attempting to open gitignore repo")

	repository, err := git.Open(storage, filesystem)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryNotExists) {
			logger.Errorf("failed to open gitignore repo: %v", err)
			return nil, err
		}

		if errors.Is(err, git.ErrRepositoryNotExists) {
			repository, err = clone(ctx, storage, filesystem)
			if err != nil {
				return nil, err
			}
			return repository, nil
		}
	}

	logger.Info("GitService initialized")
	if !options.UpdateRepository {
		return repository, nil
	}

	err = update(ctx, repository)
	if err != nil {
		return repository, err
	}

	return repository, nil
}

func clone(
	ctx context.Context,
	storage storage.Storer,
	filesystem billy.Filesystem,
) (*git.Repository, error) {
	logger := logs.CreateLogger("git.clone")
	logger.Info("cloning gitignore repository")
	repository, err := git.CloneContext(ctx, storage, filesystem, &git.CloneOptions{
		URL: GitIgnoreRepository,
	})

	if err != nil {
		message := "failed to clone gitignore repo"
		logger.Errorf("%s: %v", message, err)
		return nil, fmt.Errorf("%s: %w", message, err)
	}

	logger.Info("gitignore repository cloned successfully")
	return repository, nil
}

func update(
	ctx context.Context,
	repository *git.Repository,
) error {
	logger := logs.CreateLogger("git.update")
	logger.Info("pulling latest changes into gitignore repository")
	worktree, err := repository.Worktree()
	if err != nil {
		message := "failed to access repository worktree"
		logger.Errorf("%s: %v", message, err)
		return fmt.Errorf("%s: %w", message, err)
	}

	err = worktree.PullContext(ctx, &git.PullOptions{})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logger.Info("already up to date")
			return nil
		}

		message := "failed to pull latest changes"
		logger.Errorf("%s: %v", message, err)
		return fmt.Errorf("%s: %w", message, err)
	}

	logger.Info("pulled latest changes successfully")
	return nil
}
