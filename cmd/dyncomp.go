package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const STOP_DIRS_KEY string = "stop_dirs"
const CONFIG_PATH string = "$HOME/.config/dyncomp/"
const USAGE string = "Usage: dyncomp COMMAND_NAME\n"
const ERROR_HOME_DIR string = "Couldn't get the user home dir: %s\n"
const ERROR_CWD string = "Error while getting the current working directory: %s\n"
const ERROR_NOT_FOUND string = "Command \"%s\" not found, add it to your %s file.\n"
const ERROR_BUILDING_COMMAND string = "Error building selected command: \"%s\", error: %s\n"
const ERROR_RUNNING_COMMAND string = "Error running selected command: \"%s\", error: %s\n"
const ERROR_READING_CONFIG string = "Error reading config, error: %s\n"

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

func PullStopDirsFromConfig(writer io.Writer, configPath string) (map[string]bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(writer, ERROR_HOME_DIR, err)
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	viper.SetDefault(STOP_DIRS_KEY, []string{homeDir})

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(writer, ERROR_READING_CONFIG, err)
		return nil, err
	}

	dirsPaths := viper.GetStringSlice(STOP_DIRS_KEY)
	stopDirs := map[string]bool{}

	for _, path := range dirsPaths {
		stopDirs[path] = true
	}

	return stopDirs, nil
}

func RunCommand(writer io.Writer, cmd *cobra.Command, args []string) {
	BLUE := color.New(color.FgBlue)
	GREEN := color.New(color.FgGreen)

	if len(args) != 1 {
		fmt.Fprintf(writer, USAGE)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(writer, ERROR_CWD, err)
		return
	}

	stopDirs, err := PullStopDirsFromConfig(writer, CONFIG_PATH)
	if err != nil {
		fmt.Fprintf(writer, "%s", err)
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

	BLUE.Fprintf(writer, "Running ")
	GREEN.Fprintf(writer, "%s\n", commandString)

	command, err := BuildDynamicCommand(commandString)
	if err != nil {
		fmt.Fprintf(writer, ERROR_BUILDING_COMMAND, args[0], err)
		return
	}

	commandReader, errStdout := command.StdoutPipe()
	errReader, errStderr := command.StderrPipe()
	if errStdout != nil || errStderr != nil {
		fmt.Fprintf(writer, "Error creating the out pipe: stdout: %s, stderr: %s", errStdout, errStderr)
	}

	outScanner := bufio.NewScanner(commandReader)
	errScanner := bufio.NewScanner(errReader)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go PrintAndNotifyWaitGroup(&writer, &waitGroup, outScanner)
	go PrintAndNotifyWaitGroup(&writer, &waitGroup, errScanner)

	command.Start()
	waitGroup.Wait()

	err = command.Wait()
	if err != nil {
		blue := color.New(color.FgBlue)
		blue.Fprintf(writer, ERROR_RUNNING_COMMAND, args[0], err)
	}
}

func PrintAndNotifyWaitGroup(writer *io.Writer, waitGroup *sync.WaitGroup, scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Fprintln(*writer, scanner.Text())
	}
	waitGroup.Done()
}
