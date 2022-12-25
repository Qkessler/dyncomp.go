package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME string = "dyncomp.json"

func MergeConfigFiles(stopDirs map[string]bool, startDir string) (map[string]string, error) {
	if len(stopDirs) == 0 {
		return map[string]string{}, errors.New("Stop dirs should not be empty.")
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
			fmt.Printf("Parsed %s incorrectly, returning current config\n", path)
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
	fmt.Println("Doesn't contain a config file.")
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
