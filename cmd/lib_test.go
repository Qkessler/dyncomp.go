package cmd

import (
	"os"
	"path/filepath"
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

func TestContainsConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.Create(filepath.Join(tempDir, CONFIG_FILE_NAME))
	if err != nil {
		t.Fatalf("Error when creating the temp file for reading: %s", err)
	}

	present, _ := ContainsConfigFile(tempDir)
	if !present {
		t.Fatalf("File should be present after creating it.")
	}
	defer os.Remove(tempFile.Name())
}

func TestContainsConfigFileUnexistent(t *testing.T) {
	present, _ := ContainsConfigFile(t.TempDir())
	if present {
		t.Fatalf("File should not be present on a temporary directory.")
	}
}

func TestMergeConfigFilesEmptyStopDirs(t *testing.T) {
	_, err := MergeConfigFiles(map[string]bool{}, "")

	if err == nil {
		t.Fatalf("MergeConfigFiles should error with empty stop dirs")
	}
}

func TestMergeConfigFilesUnrelatedStopAndStart(t *testing.T) {
	_, err := MergeConfigFiles(map[string]bool{"/tmp": true}, "~/Documents")

	if err == nil {
		t.Fatalf("MergeConfigFiles should error with empty stop dirs")
	}

}

func TestMergeConfigFiles(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	tempDirConfig := filepath.Join(tempDir, CONFIG_FILE_NAME)
	subDirConfig := filepath.Join(subDir, CONFIG_FILE_NAME)

	os.Create(tempDirConfig)
	os.Mkdir(subDir, 0777)
	os.Create(subDirConfig)

	writeStringToFile(tempDirConfig, "{\"run\": \"from temp dir\"}")
	writeStringToFile(subDirConfig, "{\"run\": \"from sub dir\"}")

	config, err := MergeConfigFiles(map[string]bool{tempDir: true}, subDir)
	if err != nil {
		t.Fatalf("Error merging config files: %s", err)
	}

	if value, _ := config["run"]; value != "from sub dir" {
		t.Fatalf("Value is not correct, expecting subdir config to take place")
	}

	config, err = MergeConfigFiles(map[string]bool{tempDir: true}, tempDir)
	if value, _ := config["run"]; value != "from temp dir" {
		t.Fatalf("Value is not correct, expecting temp config to take place")
	}

	writeStringToFile(subDirConfig, "incorrect json")
	config, err = MergeConfigFiles(map[string]bool{tempDir: true}, subDir)
	if err == nil {
		t.Fatal("On incorrect subdir, we should return the empty config and error")
	}

	defer os.RemoveAll(subDir)
	defer os.Remove(subDirConfig)
	defer os.Remove(tempDirConfig)
}

func writeStringToFile(filePath string, toWrite string) error {
	fd, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	if _, err := fd.WriteString(toWrite); err != nil {
		return err
	}
	fd.Close()

	return nil
}
