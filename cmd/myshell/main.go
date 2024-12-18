package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input

		inputCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		inputCommand = strings.TrimSpace(inputCommand)

		// Exit the shell if the user types "exit 0"
		if inputCommand == "exit 0" {
			os.Exit(0)
		} else if strings.Contains(inputCommand, "echo") {
			// Print the string after "echo"
			fmt.Println(strings.TrimSpace(strings.TrimPrefix(inputCommand, "echo")))
		} else {
			fmt.Printf("%s: command not found\n", strings.TrimSpace(inputCommand))
		}
	}
}
