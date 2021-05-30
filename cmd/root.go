package cmd

import (
	"github.com/haroldadmin/getignore/internal/logs"
	"github.com/spf13/cobra"
)

var verbose bool
var extraVerbose bool

var rootCmd = &cobra.Command{
	Use:   "getignore",
	Short: "Fetch .gitignore from your terminal",
	Long: `getignore helps you fetch .gitignore files right from your terminal
	
See available commands and usage instructions using the --help flag.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		true,
		"Print info logs too",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&extraVerbose,
		"extra-verbose",
		"vv",
		false,
		"Print info and debug logs too",
	)

	logs.SetupLogs(logs.LogConfig{
		Verbose:     verbose,
		VeryVerbose: extraVerbose,
	})
}
