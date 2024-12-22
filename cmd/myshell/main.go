package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	initCommands()
	for {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		in, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("error reading from stdin: %s", err.Error())
			os.Exit(1)
		}
		inputs := strings.Split(strings.TrimSpace(in), " ")
		cmd := inputs[0]
		args := inputs[1:]
		cmdFn, ok := commands[cmd]
		if !ok {
			notFound(cmd)
		} else {
			cmdFn(args)
		}
	}
}

type (
	cmdFnc func([]string)
)

var commands = make(map[string]cmdFnc)

func registerCommand(cmd string, fn cmdFnc) {
	commands[cmd] = fn
}

func initCommands() {
	registerCommand("exit", exit)
	registerCommand("echo", echo)
	registerCommand("type", typer)
}

func notFound(cmd string) {
	fmt.Printf("%s: command not found\n", cmd)
}

func exit(args []string) {
	if len(args) == 0 {
		os.Exit(1)
	}
	if code, err := strconv.Atoi(args[0]); err == nil {
		os.Exit(code)
	}
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func typer(args []string) {
	if len(args) == 0 {
		fmt.Println(os.Stdout, "")
		fmt.Println("")
	}
	_, builtin := commands[args[0]]
	if builtin {
		fmt.Printf("%s is a shell builtin\n", args[0])
		return
	}
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		fp := filepath.Join(path, args[0])
		if _, err := os.Stat(fp); err == nil {
			fmt.Println(fp)
			return
		}
	}
	fmt.Printf("%s not found\n", args[0])
}
