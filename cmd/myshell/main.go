package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type (
	cmdFnc func([]string) string
)

var commands = make(map[string]cmdFnc)

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
func notFound(cmd string, args []string, outputFile string) {
	// Check if the command exists in the PATH, and if it does, execute it
	cmdPath, err := findCommandInPath(cmd)
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}

	// Prepare the command execution
	command := exec.Command(cmdPath, args...)

	// Redirect output to a file or Stdout
	var outputWriter io.Writer
	if outputFile != "" {
		// If the file does not exist, create it
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("error creating file: %s\n", err.Error())
			return
		}
		defer file.Close()
		// Use MultiWriter to write to both file and stdout
		outputWriter = io.MultiWriter(file, os.Stdout)
	} else {
		outputWriter = os.Stdout
	}

	// Attach output writers for Stdout and Stderr
	command.Stdout = outputWriter
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	// Run the command
	if err := command.Run(); err != nil {
		fmt.Printf("error executing command: %s\n", err.Error())
	}

	// Explicitly flush the output writer if it's a file
	if outputFile != "" {
		if f, ok := outputWriter.(*os.File); ok {
			f.Sync()
		}
	}
}

// Command functions: exit command
func exit(args []string) string {
	if len(args) == 0 {
		os.Exit(1)
	}
	if code, err := strconv.Atoi(args[0]); err == nil {
		os.Exit(code)
	}
	return ""
}

// Command functions: echo command
func echo(args []string) string {
	return strings.Join(args, " ")
}

// Command functions: type command
func type_(args []string) string {
	if len(args) == 0 {
		return "type: usage: type <command>"
	}

	command := args[0]

	// Check if the command is a shell builtin
	if command != "cat" {
		_, builtin := commands[command]
		if builtin {
			return fmt.Sprintf("%s is a shell builtin", command)
		}
	}

	// Check if the command exists in the PATH
	cmdPath, err := findCommandInPath(command)

	if err != nil {
		return fmt.Sprintf("%s: not found", command)
	} else {
		return fmt.Sprintf("%s is %s", command, cmdPath)
	}
}

// Command functions: pwd command
func pwd(args []string) string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("error getting current directory: %s", err.Error())
	}
	return dir
}

// Command functions: cd command
func cd(args []string) string {
	if len(args) == 0 {
		return "cd: usage: cd <directory>"
	}

	targetDir := args[0]

	// If the target directory is ~, change to the home directory
	if targetDir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Sprintf("error getting home directory: %s", err.Error())
		}
		targetDir = homeDir
	}

	// Change to the target directory
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Sprintf("cd: %s: No such file or directory", targetDir)
	} else {
		return ""
	}
}

// Command functions: cat command
func cat(args []string) string {
	if len(args) == 0 {
		return "cat: usage: cat <file1> <file2> ..."
	}

	output := ""

	// Read the files
	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			return fmt.Sprintf("cat: %s: No such file or directory", arg)
		}
		defer file.Close()

		// Read the file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			output += scanner.Text()
		}
	}

	return output
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
func parseInput(input string) (string, []string, string) {
	// Trim input to remove leading and trailing spaces
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", []string{}, ""
	}

	// Find the command (first word before the first space)
	spaceIndex := strings.Index(input, " ")
	var command string
	if spaceIndex == -1 {
		// If no spaces, the input is the command with no arguments
		command = input
		return command, []string{}, ""
	} else {
		command = input[:spaceIndex]
		input = input[spaceIndex+1:]
	}

	arguments := []string{}
	currentArg := ""
	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNext := false
	outputFile := ""

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
					currentArg += "\\" + string(char)
				}
			} else {
				currentArg += string(char)
			}
			escapeNext = false

		case char == '\\' && !inSingleQuotes:
			// Escape the next character
			escapeNext = true

		case char == '\'':
			// Handle single quotes
			if inDoubleQuotes {
				currentArg += string(char)
			} else {
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
		return command, []string{}, ""
	}

	// Process for redirection
	for i := 0; i < len(arguments); i++ {
		arg := arguments[i]
		if arg == ">" || arg == "1>" {
			// Ensure there's a file name after the redirection symbol
			if i+1 < len(arguments) {
				outputFile = arguments[i+1]
				// Remove the redirection symbol and the output file from arguments
				arguments = append(arguments[:i], arguments[i+2:]...)
				break
			}
		}
	}

	return command, arguments, outputFile
}

// Main function
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
		inputCommand, commandArguments, outputFile := parseInput(input)

		// Get command execution function
		execute, ok := commands[inputCommand]
		if !ok {
			notFound(inputCommand, commandArguments, outputFile)
		} else {
			result := execute(commandArguments)
			if outputFile != "" {
				// Write the output to both file and stdout
				file, err := os.Create(outputFile)
				if err != nil {
					fmt.Printf("error creating file: %s\n", err.Error())
					return
				}
				defer file.Close()
				writer := io.MultiWriter(file, os.Stdout)
				writer.Write([]byte(result))
			} else {
				fmt.Println(result)
			}
		}
	}
}
