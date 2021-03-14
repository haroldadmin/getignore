package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/haroldadmin/getignore/internal/app"
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{
	Use:   "getignore",
	Short: "Fetch .gitignore files right from the terminal",
	RunE: func(command *cobra.Command, args []string) error {
		ctx := command.Context()
		setupLogger(verbose)

		getignore, err := app.Create(ctx, app.GetIgnoreOptions{
			ShouldUpdate: noUpdate,
			Output:       outputFile,
			SearchQuery:  searchQuery,
		})

		if err != nil {
			return err
		}

		getignore.Start(ctx)
		return nil
	},
}

var (
	verbose     bool
	noUpdate    bool
	outputFile  string
	searchQuery string
)

func init() {
	flags := RootCmd.PersistentFlags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "Log extra information to the console")
	flags.BoolVar(&noUpdate, "no-update", false, "Skip pulling latest changes in gitignore repo before searching")
	flags.StringVarP(&outputFile, "output", "o", ".gitignore", "Path of .gitignore file to be written")
	flags.StringVarP(&searchQuery, "search", "s", "", "Search query to for non-interactive search")
}

func setupLogger(verbose bool) {
	level := log.ErrorLevel
	if verbose {
		level = log.DebugLevel
	}

	log.SetLevel(level)
	log.SetHandler(cli.Default)
}

func Execute(ctx context.Context) {
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
