package cmd

import (
	"log"
	"os"

	"github.com/amine7536/reverse-scan/pkg/config"
	"github.com/amine7536/reverse-scan/pkg/scanner"
	"github.com/spf13/cobra"
)

// NewRootCmd will setup and return the root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "reverse-scan",
		Short: "Reverse Scan",
		Run:   run,
	}

	rootCmd.PersistentFlags().StringP("start", "s", "", "ip range start")
	rootCmd.PersistentFlags().StringP("end", "e", "", "ip range end")
	rootCmd.PersistentFlags().StringP("output", "o", "", "csv output file")
	rootCmd.PersistentFlags().IntP("workers", "w", 8, "number of workers")

	versionCmd := NewVersionCmd()
	rootCmd.AddCommand(versionCmd)

	return rootCmd
}

func run(cmd *cobra.Command, args []string) {
	c, err := config.LoadConfig(cmd)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = scanner.Start(c)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
