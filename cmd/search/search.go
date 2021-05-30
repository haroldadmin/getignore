package search

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/haroldadmin/getignore/pkg/git"
	"github.com/haroldadmin/getignore/pkg/gitignore"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for .gitignore files interactively",
	Long: `The search command runs an interactive flow to search
for .gitignore files using their name.`,
	Args: cobra.NoArgs,
	RunE: Search,
}

var (
	repoDir      string
	updateRepo   bool
	appendToFile bool
)

var logger = log.WithField("name", "search-cmd")

func init() {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("failed to determine user's home directory: %v", err)
	}
	repoDir = filepath.Join(homeDir, ".getignore", "gitignore")

	SearchCmd.Flags().BoolVarP(
		&appendToFile,
		"append",
		"a",
		true,
		"Append to the existing .gitignore rather than overwrite it",
	)

	SearchCmd.Flags().StringVar(
		&repoDir,
		"repo-dir",
		repoDir,
		"Set custom directory for gitignore repository",
	)

	SearchCmd.Flags().BoolVar(
		&updateRepo,
		"update-repo",
		true,
		"Update the gitignore repository with upstream changes",
	)
}

func Search(cmd *cobra.Command, args []string) error {
	context := cmd.Context()
	options := git.CreateOptions{
		RepositoryDir:    repoDir,
		UpdateRepository: updateRepo,
	}

	repository, err := git.Create(context, options)
	if err != nil {
		return err
	}

	service, err := gitignore.Create(repository)
	if err != nil {
		return err
	}

	selectedFile, err := promptForQuery(context, service)
	if err != nil {
		return err
	}

	logger.Infof("selected %s", selectedFile.Name)
	workingDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to determine working directory: %v", err)
		return err
	}
	workingDirFs := osfs.New(workingDir)

	if appendToFile {
		logger.Infof("appending contents to %q", selectedFile.Name)
		err = service.Append(selectedFile, workingDirFs)
		if err != nil {
			return err
		}
		logger.Info("appended successfully")
		return nil
	}

	logger.Infof("overwriting .gitignore")
	err = service.Write(selectedFile, workingDirFs)
	if err != nil {
		return err
	}
	logger.Infof(".gitignore written successfully")

	return nil
}

func promptForQuery(
	ctx context.Context,
	service gitignore.GitIgnoreService,
) (gitignore.GitIgnoreFile, error) {
	var selectedFile gitignore.GitIgnoreFile
	for ctx.Err() == nil {
		searchPrompt := promptui.Prompt{Label: "Search"}
		query, err := searchPrompt.Run()
		if err != nil {
			return selectedFile, fmt.Errorf("prompt error: %v", err)
		}

		results, err := service.Search(query)
		if err != nil {
			if errors.Is(err, gitignore.ErrNotFound) {
				continue
			} else {
				return selectedFile, err
			}
		}

		options := make([]string, 0, len(results))
		for _, result := range results {
			options = append(options, result.Name)
		}
		options = append(options, "search again")

		selectionPrompt := promptui.Select{
			Label: "Results",
			Items: options,
		}

		index, _, err := selectionPrompt.Run()
		if err != nil {
			return selectedFile, fmt.Errorf("failed to read selection result: %v", err)
		}

		if index == len(results) {
			// User wants to retry searching
			continue
		}

		return results[index], nil
	}
	return selectedFile, errors.New("cancelled")
}
