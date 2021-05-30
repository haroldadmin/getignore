package get

import (
	"errors"

	"github.com/apex/log"
	"github.com/haroldadmin/getignore/pkg/git"
	"github.com/haroldadmin/getignore/pkg/gitignore"
	"github.com/spf13/cobra"
)

var (
	repoDir      string
	updateRepo   bool
	appendToFile bool
)

var logger = log.WithField("name", "get-cmd")

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the .gitignore file that matches the supplied name",
	Long: `get command fetches the .gitignore that best matches the supplied
name non-interactively. Exits with an error if no close match is found.`,
	Args: cobra.ExactArgs(1),
	RunE: RunGet,
}

func init() {
	GetCmd.Flags().BoolVarP(
		&appendToFile,
		"append",
		"a",
		true,
		"Append to the existing .gitignore rather than overwrite it",
	)

	GetCmd.Flags().StringVar(
		&repoDir,
		"repo-dir",
		repoDir,
		"Set custom directory for gitignore repository",
	)

	GetCmd.Flags().BoolVar(
		&updateRepo,
		"update-repo",
		true,
		"Updating the gitignore repository",
	)
}

func RunGet(cmd *cobra.Command, args []string) error {
	context := cmd.Context()
	repository, err := git.Create(context, git.CreateOptions{
		RepositoryDir:    repoDir,
		UpdateRepository: updateRepo,
	})
	if err != nil {
		return err
	}

	service, err := gitignore.Create(repository)
	if err != nil {
		return err
	}

	fileName := args[0]
	file, err := service.Get(fileName)
	if err != nil {
		if errors.Is(err, gitignore.ErrNotFound) {
			logger.Errorf("no match found for %q", fileName)
			return nil
		}

		return err
	}

	logger.Infof("selected %q", file.Name)
	return nil
}
