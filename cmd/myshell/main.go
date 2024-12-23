package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
			notFound(inputCommand, commandArguments)
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
	registerCommand("pwd", pwd)
	registerCommand("cd", cd)
}

// Function to handle command not found
func notFound(cmd string, args []string) {
	// Check if the command exists in the PATH, and if it does, execute it
	cmdPath, err := findCommandInPath(cmd)

	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
	} else {
		// Execute the command
		arg := strings.Join(args, " ")
		cmd := exec.Command(cmdPath, arg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Printf("error executing command: %s\n", err.Error())
		}
	}
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

	// Check if the command exists in the PATH
	cmdPath, err := findCommandInPath(args[0])

	if err != nil {
		fmt.Printf("%s: not found\n", args[0])
	} else {
		fmt.Printf("%s is %s\n", args[0], cmdPath)
	}
}

// Command functions: pwd command
func pwd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting current directory: %s\n", err.Error())
		return
	}
	fmt.Println(dir)
}

// Command functions: cd command
func cd(args []string) {
	if len(args) == 0 {
		fmt.Println("cd: usage: cd <directory>")
		return
	}

	if err := os.Chdir(args[0]); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
	}
}

// Helper functions
// Function to find the command in the PATH
func findCommandInPath(cmd string) (string, error) {
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		fp := filepath.Join(path, cmd)

		// Check if the file exists
		if _, err := os.Stat(fp); err == nil {
			return fp, nil
		}
	}
	return "", fmt.Errorf("command not found")
}
