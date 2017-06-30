package main

import (
	"log"

	"github.com/amine7536/reverse-scan/cmd"
)

const (
	// Version : app version
	Version = "0.2.1"
	// ProgramName : app name
	ProgramName = "Reverse-Scan"
)

func main() {

	if err := cmd.NewRootCmd(Version, ProgramName).Execute(); err != nil {
		log.Fatal(err)
	}
}
