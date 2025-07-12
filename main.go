package main

import (
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	os.Args = os.Args[1:]

	switch command {
	case "generator":
		handleGenerator()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`EOS - Framework for generation different stuff

In ancient Greek mythology and religion, Eos is the goddess and personification of the dawn,
who rose each morning from her home at the edge of the river Oceanus to deliver light and disperse the night.

Usage:
  eos <command> [flags]

Commands:
  generator    Generate different patterns and structures
  help         Show this help message

Examples:
  eos generator wrapper -n MyStruct --logging --timeout
  eos generator project -p myproject -o ./output
  eos generator repository -t User -i`)
}

func handleGenerator() {
	if len(os.Args) < 2 {
		printGeneratorUsage()
		os.Exit(1)
	}

	subCommand := os.Args[1]
	os.Args = os.Args[1:] // Remove the subcommand from args

	switch subCommand {
	case "wrapper":
		cmd.HandleWrapper()
	case "project":
		cmd.HandleProject()
	case "repository":
		cmd.HandleRepository()
	case "common":
		cmd.HandleCommon()
	// case "project-v2":
	// 	cmd.HandleProjectV2()
	case "help", "-h", "--help":
		printGeneratorUsage()
	default:
		fmt.Printf("Unknown generator command: %s\n", subCommand)
		printGeneratorUsage()
		os.Exit(1)
	}
}

func printGeneratorUsage() {
	fmt.Println(`Generator commands:

  wrapper       Generate middleware wrappers for Go structs
  project       Generate a new Go project structure
  repository    Generate repository pattern for structs
  common        Generate common files

Examples:
  eos generator wrapper -n MyStruct --logging --timeout
  eos generator project -p myproject -o ./output
  eos generator repository -t User -i`)
}
