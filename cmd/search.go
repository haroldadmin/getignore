package cmd

import (
	"path/filepath"

	"github.com/apex/log"
	"github.com/haroldadmin/getignore/pkg/git"
	"github.com/haroldadmin/getignore/pkg/gitignore"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for .gitignore files interactively",
	Long: `The search command runs an interactive flow to search
for .gitignore files using their name.`,
	Args: cobra.ExactArgs(1),
	RunE: Search,
}

var (
	repoDir    string
	updateRepo bool
)

func init() {
	rootCmd.AddCommand(searchCmd)

	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("failed to determine user's home directory: %v", err)
	}
	repoDir = filepath.Join(homeDir, ".getignore", "gitignore")

	searchCmd.Flags().StringVar(
		&repoDir,
		"repo-dir",
		repoDir,
		"Set custom directory for gitignore repository",
	)

	searchCmd.Flags().BoolVar(
		&updateRepo,
		"update-repo",
		true,
		"Updating the gitignore repository",
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

	searchQuery := args[0]
	service.Search(searchQuery)

	return nil
}
