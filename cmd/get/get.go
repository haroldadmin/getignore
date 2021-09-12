package get

import (
	"errors"
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/haroldadmin/getignore/internal/logs"
	"github.com/haroldadmin/getignore/pkg/git"
	"github.com/haroldadmin/getignore/pkg/gitignore"
	"github.com/spf13/cobra"
)

var (
	repoDir      string
	updateRepo   bool
	appendToFile bool
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the .gitignore file with the exact given name",
	Long: `Non-interactively fetches the .gitignore file with the exact name. 
Exits with an error if no match is found.

Use this command if you're sure of the exact name of the .gitignore file
you're looking for.`,
	Args: cobra.ExactArgs(1),
	RunE: RunGet,
}

func init() {
	GetCmd.Flags().BoolVarP(
		&appendToFile,
		"append",
		"a",
		true,
		`Append to the existing .gitignore rather than overwrite it. 
Creates a new .gitignore file if it doesn't exist.`,
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
	logger := logs.CreateLogger("cmd.get")
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
		logger.Infof("appending contents to %q", file.Name)
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
