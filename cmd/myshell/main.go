package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {

	for {
		// Print the shell prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		inputCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		inputCommand = strings.TrimSpace(inputCommand)

		// Handle the command
		output := handleCommand(inputCommand)
		fmt.Println(output)
	}
}

func handleCommand(inputCommand string) string {
	// List of built-in commands
	builtInCommands := []string{"echo", "exit", "type"}

	// Exit the shell if the user types "exit 0"
	if inputCommand == "exit 0" {
		os.Exit(0)
	} else if strings.HasPrefix(inputCommand, "echo") {
		// Print the string after "echo"
		return strings.TrimSpace(strings.TrimPrefix(inputCommand, "echo"))
	} else if strings.HasPrefix(inputCommand, "type") {
		// Print the type of the command
		command := strings.TrimSpace(strings.TrimPrefix(inputCommand, "type"))
		if slices.Contains(builtInCommands, command) {
			return fmt.Sprintf("%s is a shell builtin\n", command)
		} else {
			return fmt.Sprintf("%s: not found\n", command)
		}
	} else {
		return fmt.Sprintf("%s: command not found\n", strings.TrimSpace(inputCommand))
	}
	return ""
}
