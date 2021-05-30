package get

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/go-git/go-billy/v5/osfs"
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
		"Update the gitignore repository with upstream changes",
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

	workingDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to determine working directory: %v", err)
		return err
	}
	workingDirFs := osfs.New(workingDir)

	if appendToFile {
		logger.Infof("appending contents to %q")
		err = service.Append(file, workingDirFs)
		if err != nil {
			return err
		}
		logger.Info("appended successfully")
		return nil
	}

	logger.Infof("overwriting .gitignore")
	err = service.Write(file, workingDirFs)
	if err != nil {
		return err
	}
	logger.Infof(".gitignore written successfully")

	return nil
}
