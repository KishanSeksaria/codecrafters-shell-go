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
		inputCommand, commandArguments := parseInput(input)

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
	registerCommand("type", type_)
	registerCommand("pwd", pwd)
	registerCommand("cd", cd)
	registerCommand("cat", cat)
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
func type_(args []string) {
	if len(args) == 0 {
		fmt.Println("type: usage: type <command>")
		return
	}

	command := args[0]

	// Check if the command is a shell builtin
	if command != "cat" {
		_, builtin := commands[command]
		if builtin {
			fmt.Printf("%s is a shell builtin\n", command)
			return
		}
	}

	// Check if the command exists in the PATH
	cmdPath, err := findCommandInPath(command)

	if err != nil {
		fmt.Printf("%s: not found\n", command)
	} else {
		fmt.Printf("%s is %s\n", command, cmdPath)
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

	targetDir := args[0]

	// If the target directory is ~, change to the home directory
	if targetDir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("error getting home directory: %s\n", err.Error())
			return
		}
		targetDir = homeDir
	}

	// Change to the target directory
	if err := os.Chdir(targetDir); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", targetDir)
	}
}

// Command functions: cat command
func cat(args []string) {
	if len(args) == 0 {
		fmt.Println("cat: usage: cat <file>")
		return
	}

	output := ""

	// Read the files
	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			fmt.Printf("cat: %s: No such file or directory\n", arg)
			return
		}
		defer file.Close()

		// Read the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			output += scanner.Text()
		}
	}

	fmt.Println(output)
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

// Function to parse input into command and arguments
func parseInput(input string) (string, []string) {
	// Trim input to remove leading and trailing spaces
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", []string{}
	}

	// Find the command (first word before the first space)
	spaceIndex := strings.Index(input, " ")
	var command string
	if spaceIndex == -1 {
		// If no spaces, the input is the command with no arguments
		command = input
		return command, []string{}
	} else {
		command = input[:spaceIndex]
		input = input[spaceIndex+1:]
	}

	arguments := []string{}
	currentArg := ""
	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNext := false

	// Parse each character
	for _, char := range input {
		switch {
		case escapeNext:
			// Handle escaped characters
			if inDoubleQuotes {
				// Allow escaping only for specific characters inside double quotes
				if char == '$' || char == '`' || char == '"' || char == '\\' || char == '\n' {
					currentArg += string(char)
				} else {
					// Treat the backslash as a literal character if the char isn't escapable
					currentArg += "\\" + string(char)
				}
			} else {
				// Outside quotes, escape any character
				currentArg += string(char)
			}
			escapeNext = false

		case char == '\\' && !inSingleQuotes:
			// Escape the next character
			escapeNext = true

		case char == '\'':
			// Handle single quotes
			if inDoubleQuotes {
				// Add single quotes literally when inside double quotes
				currentArg += string(char)
			} else {
				// Toggle single-quote state outside of double quotes
				inSingleQuotes = !inSingleQuotes
			}

		case char == '"' && !inSingleQuotes:
			// Toggle double-quote state
			inDoubleQuotes = !inDoubleQuotes

		case char == ' ' && !inSingleQuotes && !inDoubleQuotes:
			// Space ends the current argument if not in quotes
			if len(currentArg) > 0 {
				arguments = append(arguments, currentArg)
				currentArg = ""
			}

		default:
			// Add the character to the current argument
			currentArg += string(char)
		}
	}

	// Add the last argument if it's not empty
	if len(currentArg) > 0 {
		arguments = append(arguments, currentArg)
	}

	// Handle mismatched quotes
	if inSingleQuotes || inDoubleQuotes {
		return command, []string{}
	}

	return command, arguments
}
