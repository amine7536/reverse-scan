package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.pixelfactory.io/pkg/version"
)

// NewVersionCmd will setup and return the version command
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s\n", version.REVISION)
			fmt.Printf("build-date: %s\n", version.BUILDDATE)
		},
	}

	return versionCmd
}
