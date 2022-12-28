package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

const FOR_TESTS_KEY string = "for_tests"
const UNEXISTENT_KEY string = "unexistent_key"
const EMPTY_KEY string = "empty_key"
const SHOULD_ERROR_KEY string = "should_error_key"

func TestDyncompCommandCreated(t *testing.T) {
	command := CreateCommand()

	if command.Use == "" || command.Example == "" || command.Run == nil {
		t.Fatalf("Command has been created incorrectly, core fields not defined: %+v", command)
	}
}

func TestRunCommandWithoutArgs(t *testing.T) {
	var outputNoArgs bytes.Buffer
	var outputTwoArgs bytes.Buffer

	RunCommand(&outputNoArgs, CreateCommand(), []string{})
	RunCommand(&outputTwoArgs, CreateCommand(), []string{"first", "second"})

	if outputNoArgs.String() != USAGE || outputTwoArgs.String() != USAGE {
		fmt.Println(outputNoArgs.String(), outputTwoArgs.String())
		t.Fatalf("Should print usage on no args, or not 1.")
	}
}

func TestRunCommandErrorMergeConflictFiles(t *testing.T) {
	_, fileName, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(fileName)
	os.Chdir(t.TempDir())
	defer os.Chdir(currentDir)

	var output bytes.Buffer

	RunCommand(&output, CreateCommand(), []string{FOR_TESTS_KEY})

	if output.String() != ERROR_START_DIR_CONTAINED {
		t.Fatalf("With directory changed, should have incorrect start dir.")
	}
}

func TestRunCommandUnexistentConfigKey(t *testing.T) {
	var output bytes.Buffer

	RunCommand(&output, CreateCommand(), []string{UNEXISTENT_KEY})

	if output.String() != fmt.Sprintf(ERROR_NOT_FOUND, UNEXISTENT_KEY, CONFIG_FILE_NAME) {
		t.Fatalf("Should not find the %s command", UNEXISTENT_KEY)
	}
}

func TestRunCommandIncorrectlyFormattedCommandString(t *testing.T) {
	var output bytes.Buffer

	RunCommand(&output, CreateCommand(), []string{EMPTY_KEY})

	if output.String() != fmt.Sprintf(ERROR_BUILDING_COMMAND, EMPTY_KEY, "invalid syntax") {
		t.Fatalf("Should error building command: %s", EMPTY_KEY)
	}
}

func TestRunCommandErrorRunningCommand(t *testing.T) {
	var output bytes.Buffer

	RunCommand(&output, CreateCommand(), []string{SHOULD_ERROR_KEY})

	if output.String() != fmt.Sprintf(ERROR_RUNNING_COMMAND, SHOULD_ERROR_KEY,
		`exec: "unexistent_dyncomp_command": executable file not found in $PATH
`) {
		t.Fatalf("Should error running command: %s", SHOULD_ERROR_KEY)
	}
}

func TestRunCommandHappyPathWithForTestsKey(t *testing.T) {
	var output bytes.Buffer

	RunCommand(&output, CreateCommand(), []string{FOR_TESTS_KEY})

	if output.String() != "EXPECTED TEST OUTPUT\n" {
		t.Fatalf("%s should be present on dyncomp.json file, and echo 'EXPECTED TESTS OUTPUT'", FOR_TESTS_KEY)
	}
}

func TestPullStopDirsFromConfigEmptyShouldHaveHome(t *testing.T) {
	var output bytes.Buffer
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "config.json")
	tempFile, err := os.Create(path)
	if err != nil {
		t.Fatalf("Shouldn't error when creating config file: %s", err)
	}
	defer os.Remove(path)

	tempFile.WriteString(`{}`)

	viper.Reset()
	stopDirs, err := PullStopDirsFromConfig(&output, tempDir)
	if err != nil {
		t.Fatalf("Shouldn't error when files created correctly: %s", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir should be defined correctly.")
	}

	if !reflect.DeepEqual(stopDirs, map[string]bool{ homeDir: true }) {
		t.Fatalf("With empty config, we should have HOME as stop dir.")
	}
}

func TestPullStopDirsFromConfigSetShouldBeCorrect(t *testing.T) {
	var output bytes.Buffer
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "config.json")
	tempFile, err := os.Create(path)
	if err != nil {
		t.Fatalf("Shouldn't error when creating config file: %s", err)
	}
	defer os.Remove(path)

	tempFile.WriteString(fmt.Sprintf(`{"%s": ["test path"]}`, STOP_DIRS_KEY))

	viper.Reset()
	stopDirs, err := PullStopDirsFromConfig(&output, tempDir)
	if err != nil {
		t.Fatalf("Shouldn't error when files created correctly: %s", err)
	}

	if !reflect.DeepEqual(stopDirs, map[string]bool{ "test path": true }) {
		t.Fatalf("With set config, it should have the selected one.")
	}
}
