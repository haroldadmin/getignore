package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/haroldadmin/getignore/app"
	"github.com/spf13/cobra"
)

var RootCmd *cobra.Command = &cobra.Command{
	Use:   "getignore",
	Short: "Fetch .gitignore files right from the terminal",
	RunE: func(command *cobra.Command, args []string) error {
		ctx := command.Context()

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
	flags.StringVarP(&searchQuery, "search", "s", "", "Search query for matching gitignore files")
}

func Execute(ctx context.Context) {
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
