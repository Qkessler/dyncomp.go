package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const USAGE string = "Usage: dyncomp COMMAND_NAME\n"
const ERROR_HOME_DIR string = "Couldn't get the user home dir: %s\n"
const ERROR_CWD string = "Error while getting the current working directory: %s\n"
const ERROR_NOT_FOUND string = "Command \"%s\" not found, add it to your %s file.\n"
const ERROR_BUILDING_COMMAND string = "Error building selected command: \"%s\", error: %s\n"
const ERROR_RUNNING_COMMAND string = "Error running selected command: \"%s\", error: %s\n"

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
		Run: run,
	}
}

func Execute() {
	error := CreateCommand().Execute()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	RunCommand(os.Stdout, cmd, args)
}

func RunCommand(writer io.Writer, cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(writer, USAGE)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(writer, ERROR_HOME_DIR, err)
		return
	}
	stopDirs := map[string]bool{
		homeDir: true,
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(writer, ERROR_CWD, err)
		return
	}
	configFiles, err := MergeConfigFiles(stopDirs, cwd)
	if err != nil {
		fmt.Fprintf(writer, "%s", err)
		return
	}

	commandString, present := configFiles[args[0]]
	if !present {
		fmt.Fprintf(writer, ERROR_NOT_FOUND, args[0], CONFIG_FILE_NAME)
		return
	}
	command, err := BuildDynamicCommand(commandString)
	if err != nil {
		fmt.Fprintf(writer, ERROR_BUILDING_COMMAND, args[0], err)
		return
	}

	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Fprintf(writer, ERROR_RUNNING_COMMAND, args[0], err)
	}

	fmt.Fprintln(writer, string(output[:]))
}
