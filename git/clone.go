package git

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	goGit "github.com/go-git/go-git/v5"
)

const gitIgnoreRepo = "https://github.com/github/gitignore.git"

func Clone(skipUpdate bool) (string, error) {
	repoDir, err := getRepoDir()
	if err != nil {
		return "", err
	}

	repo, err := initRepo(repoDir)
	if err != nil {
		return "", err
	}

	if !skipUpdate {
		err = updateRepo(repo)
		if err != nil {
			return "", err
		}
	}

	log.Printf("Repository initialized at %v", repoDir)
	return repoDir, nil
}

func getRepoDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not access home directory: %v", err)
	}

	repoDir := filepath.Join(home, ".getignore", "gitignore")
	if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("Failed to access %v: %v", repoDir, err)
	}

	return repoDir, nil
}

func initRepo(repoDir string) (*goGit.Repository, error) {
	log.Printf("Attempting to open %v as a repository", repoDir)
	repo, err := goGit.PlainOpen(repoDir)

	if err != nil {
		if err != goGit.ErrRepositoryNotExists {
			return nil, fmt.Errorf("%v is not a valid git repository: %v", repoDir, err)
		}

		log.Printf("Repository does not exist, attempting to clone")
		repo, err = goGit.PlainClone(repoDir, false, &goGit.CloneOptions{
			URL:      gitIgnoreRepo,
			Depth:    1,
			Tags:     goGit.NoTags,
			Progress: os.Stdout,
		})

		if err != nil {
			return nil, fmt.Errorf("Failed to clone repository: %v", err)
		}
	}

	return repo, nil
}

func updateRepo(repo *goGit.Repository) error {
	log.Printf("Attempting to update the repository with latest changes")
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Failed to get working tree of repo: %v", err)
	}

	err = workTree.Pull(&goGit.PullOptions{RemoteName: "origin"})
	if err != nil && err != goGit.NoErrAlreadyUpToDate {
		return fmt.Errorf("Failed to pull latest changes: %v", err)
	}

	return nil
}
