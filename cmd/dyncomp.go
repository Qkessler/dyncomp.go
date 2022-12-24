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
			fmt.Println("Testing whether the command runs.")
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
