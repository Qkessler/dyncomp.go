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
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Walk the directory tree upwards and get all dyncomp.json files
			// up until the home directory.
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Couldn't get the user home dir: %s", err)
				return
			}
			stopDirs := map[string]bool{
				homeDir: true,
			}
			configFiles, err := PullConfigFiles(stopDirs)

			fmt.Println(configFiles)
			// TODO: Compile all the commands into a map string
			// TODO: Understand how to build the commands async with Go. Go Coroutines?
		},
	}
}

func Execute() {
	error := CreateCommand().Execute()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}
