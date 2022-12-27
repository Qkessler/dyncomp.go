package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	CommandName string
}

func CreateCommand() *cobra.Command {
	return &cobra.Command{
		Use: "dyncomp COMMAND_NAME",
		Example: `These commands assume that you have defined the "run" and the "test"
keys on your dyncomp.json configuration file.

- dyncomp run
- dyncomp test`,
		Run: RunCommand,
	}
}

func Execute() {
	error := CreateCommand().Execute()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}

func RunCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Printf("Usage: dyncomp COMMAND_NAME\n")
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Couldn't get the user home dir: %s\n", err)
		return
	}
	stopDirs := map[string]bool{
		homeDir: true,
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting the current working directory: %s\n", err)
		return
	}
	configFiles, err := MergeConfigFiles(stopDirs, cwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	commandName, present := configFiles[args[0]]
	if !present {
		fmt.Printf("Command \"%s\" not found, add it to your %s file.\n", args[0], CONFIG_FILE_NAME)
		return
	}
	command, err := BuildDynamicCommand(commandName)
	if err != nil {
		fmt.Printf("Error building selected command: \"%s\", error: %s\n", args[0], err)
		return
	}

	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running selected command: \"%s\", error: %s\n", args[0], err)
	}

	fmt.Println(string(output[:]))
}
