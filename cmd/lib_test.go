package cmd

import (
	"os"
	"testing"
)

func TestParseConfigFileNonExistent(t *testing.T) {
	_, err := ParseConfigFile("unexistent-file.json")
	if err == nil {
		t.Fatal("ParseConfigFile should fail, unexsistent file.")
	}
}

func TestParseConfigFileEmpty(t *testing.T) {
	tempFile, err := os.CreateTemp("", "dyncomp.json")
	tempFile.WriteString("{}")
	if err != nil {
		t.Fatalf("Error when creating the temp file for reading: %s", err)
	}

	configMap, err := ParseConfigFile(tempFile.Name())
	if err != nil {
		t.Fatalf("ParseConfigFile shouldn't fail here: %s", err)
	}

	if len(configMap) != 0 {
		t.Fatalf("ConfigMap should be empty")
	}
	defer os.Remove(tempFile.Name())
}

func TestParseConfigFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "dyncomp.json")
	tempFile.WriteString("{\"run\": \"go run main.go\"}")
	if err != nil {
		t.Fatalf("Error when creating the temp file for reading: %s", err)
	}

	configMap, err := ParseConfigFile(tempFile.Name())
	if err != nil {
		t.Fatalf("ParseConfigFile shouldn't fail here: %s", err)
	}

	if configMap["run"] == "" {
		t.Fatalf("Run key should have value after parsing.")
	}
	defer os.Remove(tempFile.Name())
}
