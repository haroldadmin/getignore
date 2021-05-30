package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appendToFile bool

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the .gitignore file that matches the supplied name",
	Long: `get command fetches the .gitignore that best matches the supplied
name non-interactively. Exits with an error if no close match is found.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&appendToFile, "append", "a", true, "Append to the existing .gitignore rather than overwrite it")
}
