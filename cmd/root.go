package cmd

import (
	"log"

	"bitbucket.org/aminebenseddik/reverse-scan/conf"
	"bitbucket.org/aminebenseddik/reverse-scan/scanner"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "reverse-scan",
	Short: "Revere Lookup",
	Run:   run,
}

var version string
var progName string

// NewRootCmd will setup and return the root command
func NewRootCmd(v string, p string) *cobra.Command {
	// Set Version and ProgramName
	version = v
	progName = p

	rootCmd.PersistentFlags().StringP("start", "s", "", "Range Start")
	rootCmd.PersistentFlags().StringP("end", "e", "", "Range End")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output File")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(cmd)
	if err != nil {
		log.Fatal(err)
	}

	scanner.Start(config)
}
