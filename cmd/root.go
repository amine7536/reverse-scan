// Package cmd provides the command-line interface for reverse-scan
package cmd

import (
	"log"
	"os"

	"github.com/amine7536/reverse-scan/pkg/config"
	"github.com/amine7536/reverse-scan/pkg/scanner"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "reverse-scan",
	Short: "Reverse Scan",
	Run:   run,
}

var version string

// NewRootCmd will setup and return the root command
func NewRootCmd(v string) *cobra.Command {
	// Set Version and ProgramName
	version = v

	rootCmd.PersistentFlags().StringP("start", "s", "", "ip range start")
	rootCmd.PersistentFlags().StringP("end", "e", "", "ip range end")
	rootCmd.PersistentFlags().StringP("output", "o", "", "csv output file")
	rootCmd.PersistentFlags().IntP("workers", "w", 8, "number of workers")

	return &rootCmd
}

func run(cmd *cobra.Command, _ []string) {
	c, err := config.LoadConfig(cmd)
	if err != nil {
		log.Fatal(err)
	}

	scanner.Start(c)
	os.Exit(0)
}
