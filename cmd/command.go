package cmd

import (
	"os/exec"

	"github.com/cosiner/argv"
)

func BuildDynamicCommand(command string) (*exec.Cmd, error) {
	commands, err := argv.Argv(command, func(backquoted string) (string, error) {
		return backquoted, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Implement piping. Right now, we take the first command if
	// the command the user is trying to pass is piped into another one.
	// Example: With `ls echo / | wc -l`, we'll only run `ls echo /`
	return exec.Command(commands[0][0], commands[0][1:]...), nil
}
