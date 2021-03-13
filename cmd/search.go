package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var SearchCommand *cobra.Command = &cobra.Command{
	Use:   "search",
	Short: "Search for gitignore files non-interactively",
	RunE: func(command *cobra.Command, args []string) error {
		if len(args) == 0 || args[0] == "" {
			return errors.New("Empty search query")
		}

		query := args[0]
		fmt.Println(query)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(SearchCommand)
}
