package main

import (
	"log"

	"github.com/kj455/intelli-cli/cmd"
	"github.com/kj455/intelli-cli/secret"
)

func main() {
	err := secret.SetupSecretIfNeeded()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
