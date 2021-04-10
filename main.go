package main

import (
	"fmt"
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
	welcomeMessage := fmt.Sprintf("%s\n", `
 _______  _______        ___      _______  _______  ______   _______  _______  _______  _______ 
|       ||       |      |   |    |       ||   _   ||      | |       ||       ||       ||       |
|    ___||   _   | ____ |   |    |   _   ||  |_|  ||  _    ||_     _||    ___||  _____||_     _|
|   | __ |  | |  ||____||   |    |  | |  ||       || | |   |  |   |  |   |___ | |_____   |   |  
|   ||  ||  |_|  |      |   |___ |  |_|  ||       || |_|   |  |   |  |    ___||_____  |  |   |  
|   |_| ||       |      |       ||       ||   _   ||       |  |   |  |   |___  _____| |  |   |  
|_______||_______|      |_______||_______||__| |__||______|   |___|  |_______||_______|  |___|  
`)
	fmt.Println(welcomeMessage)

	args := os.Args[1:]
	// Validate commands
	if ok := validateCommand(args[0]); !ok {
		log.Fatal("Command not valid. this are the allowed commands: ", allowedCommands)
	}

	cmd := commands.NewCmd(args[0])
	args = args[1:]
	// Validate arguments and execute command
	if argsErr := cmd.Execute(args...); argsErr != nil {
		log.Fatal("error runnning command", argsErr.Error())
	}
}