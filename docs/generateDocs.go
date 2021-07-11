package main

import (
	"log"

	"github.com/rspace-os/rspace-cli/cmd"

	"github.com/spf13/cobra/doc"
)

func main() {
	operatorCommand := cmd.NewOperatorCommand()

	err := doc.GenMarkdownTree(operatorCommand, "./generated")
	if err != nil {
		log.Fatal(err)
	}
}
