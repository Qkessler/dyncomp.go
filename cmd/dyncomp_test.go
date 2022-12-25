package cmd

import (
	"testing"
)

func TestDyncompCommandCreated(t *testing.T) {
	command := CreateCommand()

	if command.Use == "" || command.Example == "" || command.Run == nil {
		t.Fatalf("Command has been created incorrectly, core fields not defined: %+v", command)
	}
}
