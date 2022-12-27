package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME string = "dyncomp.json"
const ERROR_START_DIR_CONTAINED string = "Start dir should be contained whithin the stop dirs\n."
const ERROR_STOP_DIRS_EMPTY string = "Stop dirs should not be empty.\n"

func MergeConfigFiles(stopDirs map[string]bool, startDir string) (map[string]string, error) {
	if len(stopDirs) == 0 {
		return map[string]string{}, errors.New(ERROR_STOP_DIRS_EMPTY)
	}

	if !any(stopDirs, func(stopDir string) bool {
		return filepath.HasPrefix(startDir, stopDir)
	}) {
		return map[string]string{},
			errors.New(ERROR_START_DIR_CONTAINED)
	}

	configPathsChannel := make(chan string)
	go func() {
		currentDir := startDir
		for {
			present, configPath := ContainsConfigFile(currentDir)
			if present {
				configPathsChannel <- configPath
			}

			if stopDirs[currentDir] {
				break
			}

			currentDir = filepath.Dir(currentDir)
		}

		close(configPathsChannel)
	}()

	config := map[string]string{}
	for path := range configPathsChannel {
		pathConfig, err := ParseConfigFile(path)
		if err != nil {
			return config, err
		}

		for key, value := range pathConfig {
			if _, containsValue := config[key]; !containsValue {
				config[key] = value
			}
		}
	}

	return config, nil
}

func ContainsConfigFile(directory string) (bool, string) {
	configPath := filepath.Join(directory, CONFIG_FILE_NAME)
	if _, err := os.Stat(configPath); err == nil {
		return true, configPath
	}
	return false, ""
}

func ParseConfigFile(filePath string) (map[string]string, error) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configFileMap map[string]string
	error := json.Unmarshal(contents, &configFileMap)
	return configFileMap, error
}

func any(collection map[string]bool, f func(string) bool) bool {
	for key := range collection {
		if f(key) {
			return true
		}
	}
	return false
}
