// Package main is the entry point for the reverse-scan application
package main

import (
	"log"

	"github.com/amine7536/reverse-scan/cmd"
)

var (
	// Version : app version (injected by goreleaser at build time)
	Version = "dev"
)

func main() {

	if err := cmd.NewRootCmd(Version).Execute(); err != nil {
		log.Fatal(err)
	}
}
