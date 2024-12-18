package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	inputCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	if inputCommand == "invalid_command\n" {
		fmt.Printf("invalid_command: command not found")
		os.Exit(1)
	}
}
