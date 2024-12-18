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

	for {
		inputCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		fmt.Printf("%s: command not found\n", strings.TrimSpace(inputCommand))
	}
}
