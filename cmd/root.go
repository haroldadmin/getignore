package cmd

import (
	"github.com/haroldadmin/getignore/cmd/get"
	"github.com/haroldadmin/getignore/cmd/search"
	"github.com/haroldadmin/getignore/internal/logs"
	"github.com/spf13/cobra"
)

var verbose bool
var extraVerbose bool

var RootCmd = &cobra.Command{
	Use:   "getignore",
	Short: "Fetch .gitignore from your terminal",
	Long: `getignore helps you fetch .gitignore files right from your terminal
	
See available commands and usage instructions using the --help flag.`,
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(
		&extraVerbose,
		"extra-verbose",
		"V",
		false,
		"Print info and debug logs too",
	)

	RootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		true,
		"Print info logs too",
	)

	RootCmd.AddCommand(get.GetCmd)
	RootCmd.AddCommand(search.SearchCmd)

	logs.SetupLogs(logs.LogConfig{
		Verbose:     verbose,
		VeryVerbose: extraVerbose,
	})
}
