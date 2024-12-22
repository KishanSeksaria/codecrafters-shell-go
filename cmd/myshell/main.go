package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type (
	cmdFnc func([]string)
)

var commands = make(map[string]cmdFnc)

func main() {
	// Initialize commands
	initCommands()

	// Main loop
	for {
		fmt.Print("$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("error reading from stdin: %s", err.Error())
			os.Exit(1)
		}

		// Parse input
		inputs := strings.Split(strings.TrimSpace(input), " ")
		inputCommand := inputs[0]
		commandArguments := inputs[1:]

		// Get command execution function
		execute, ok := commands[inputCommand]
		if !ok {
			notFound(inputCommand)
		} else {
			execute(commandArguments)
		}
	}
}

// Function to register a command
func registerCommand(cmd string, fn cmdFnc) {
	commands[cmd] = fn
}

// Function to initialize commands
func initCommands() {
	registerCommand("exit", exit)
	registerCommand("echo", echo)
	registerCommand("type", typer)
}

// Function to handle command not found
func notFound(cmd string) {
	fmt.Printf("%s: command not found\n", cmd)
}

// Command functions: exit command
func exit(args []string) {
	if len(args) == 0 {
		os.Exit(1)
	}
	if code, err := strconv.Atoi(args[0]); err == nil {
		os.Exit(code)
	}
}

// Command functions: echo command
func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// Command functions: type command
func typer(args []string) {
	if len(args) == 0 {
		fmt.Println("type: usage: type <command>")
		return
	}

	// Check if the command is a shell builtin
	_, builtin := commands[args[0]]
	if builtin {
		fmt.Printf("%s is a shell builtin\n", args[0])
		return
	}

	// Check if the command is in the PATH
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		fp := filepath.Join(path, args[0])

		// Check if the file exists
		if _, err := os.Stat(fp); err == nil {
			fmt.Println(fp)
			return
		}
	}

	fmt.Printf("%s: not found\n", args[0])
}
