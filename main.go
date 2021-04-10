package main

import (
	"go-loadtest/commands"
	"log"
	"os"
)

/*
	OpenSource. No license required.
	This cli is to create a load test for testing purpose
	Use it with responsability.
    @Version 1.0.0 - 2021
*/

var (
	allowedCommands = []string{"start"}
)

func validateCommand(cmd string) bool {
	for _, command := range allowedCommands {
		if command == cmd {
			return true
		}
	}
	return false
}

func main() {
	args := os.Args[1:]
	// Validate commands
	if ok := validateCommand(args[0]); !ok {
		log.Fatal("Command not valid. this are the allowed commands: ", allowedCommands)
	}

	cmd := commands.NewCmd(args[0])
	args = args[1:]

	// Validate arguments
	if argsErr := cmd.Validate(args...); argsErr != nil {
		log.Fatal("arguments not valid.", argsErr.Error())
	}

	// Execute Command
	cmd.Execute()
}