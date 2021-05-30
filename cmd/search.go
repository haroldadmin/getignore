package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for .gitignore files interactively",
	Long: `The search command runs an interactive flow to search
for .gitignore files using their name.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("search called")
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
