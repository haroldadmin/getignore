package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	goGit "github.com/go-git/go-git/v5"
)

const gitIgnoreRepo = "https://github.com/github/gitignore.git"

func Clone(ctx context.Context, skipUpdate bool) (string, error) {
	repoDir, err := getRepoDir()
	if err != nil {
		return "", err
	}

	repo, err := initRepo(ctx, repoDir)
	if err != nil {
		return "", err
	}

	if !skipUpdate {
		err = updateRepo(ctx, repo)
		if err != nil {
			return "", err
		}
	}

	log.Debugf("Repository initialized at %v", repoDir)
	return repoDir, nil
}

func getRepoDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not access home directory: %v", err)
	}

	repoDir := filepath.Join(home, ".getignore", "gitignore")
	if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to access %v: %v", repoDir, err)
	}

	return repoDir, nil
}

func initRepo(ctx context.Context, repoDir string) (*goGit.Repository, error) {
	log.Debugf("Attempting to open %v as a repository", repoDir)
	repo, err := goGit.PlainOpen(repoDir)

	if err != nil {
		if err != goGit.ErrRepositoryNotExists {
			return nil, fmt.Errorf("%v is not a valid git repository: %v", repoDir, err)
		}

		log.Debugf("Repository does not exist, attempting to clone")
		repo, err = goGit.PlainCloneContext(ctx, repoDir, false, &goGit.CloneOptions{
			URL:      gitIgnoreRepo,
			Depth:    1,
			Tags:     goGit.NoTags,
			Progress: os.Stdout,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to clone repository: %v", err)
		}
	}

	return repo, nil
}

func updateRepo(ctx context.Context, repo *goGit.Repository) error {
	log.Debugf("Attempting to update the repository with latest changes")
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get working tree of repo: %v", err)
	}

	err = workTree.PullContext(ctx, &goGit.PullOptions{RemoteName: "origin"})
	if err != nil && err != goGit.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull latest changes: %v", err)
	}

	return nil
}
