package main

import (
	"log"

	"bitbucket.org/aminebenseddik/reverse-scan/cmd"
)

const (
	// Version : app version
	Version = "0.1.0"
	// ProgramName : app name
	ProgramName = "Reverse-Scan"
)

func main() {

	if err := cmd.NewRootCmd(Version, ProgramName).Execute(); err != nil {
		log.Fatal(err)
	}
}
