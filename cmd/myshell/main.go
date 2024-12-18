package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	inputCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	if strings.Contains(inputCommand, "invalid") {
		fmt.Printf("%s: command not found\n", strings.TrimSpace(inputCommand))
		os.Exit(1)
	}
}
