package cmd

import (
	"testing"
)

func TestBuildDynamicCommand(t *testing.T) {
	command, err := BuildDynamicCommand("ls   `echo /`   |  wc  -l")
	if err != nil {
		t.Fatal("Should be able to parse this command: ", err)
	}

	if command.String() != "/bin/ls echo /" {
		t.Fatal("Command should be properly parsed and augmented.")
	}
}

func TestBuildDynamicCommandEmpty(t *testing.T) {
	_, err := BuildDynamicCommand("")
	if err == nil {
		t.Fatal("BuildDynamicCommand should error if no command is passed.")
	}
}
