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

		if strings.TrimSpace(inputCommand) == "exit 0" {
			// exit with status 0
			os.Exit(0)
		}

		fmt.Printf("%s: command not found\n", strings.TrimSpace(inputCommand))
	}
}
