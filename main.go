package main

import (
	"log"

	"github.com/amine7536/reverse-scan/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.pixelfactory.io/pkg/version"
)

func initConfig() {
	viper.Set("revision", version.REVISION)
}

func main() {
	cobra.OnInitialize(initConfig)

	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
