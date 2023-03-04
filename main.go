package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kj455/intelli-cli/cmd"
)

func main() {
	loadEnv()

	cmd.Execute()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}
}
