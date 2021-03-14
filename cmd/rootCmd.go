package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command = &cobra.Command{
	Use:   "getignore",
	Short: "Fetch .gitignore files right from the terminal",
	RunE: func(command *cobra.Command, args []string) error {
		fmt.Println("Root command")
		return nil
	},
}

var (
	verbose    bool
	skipUpdate bool
	outputFile string
)

func init() {
	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "Log extra information to the console")
	flags.BoolVar(&skipUpdate, "no-update", false, "Skip pulling latest changes in gitignore repo before searching")
	flags.StringVarP(&outputFile, "output", "o", ".gitignore", "Path of .gitignore file to be written")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
