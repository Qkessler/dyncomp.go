package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: Think about where to stop, since it shouldn't work with the HOME directory
// when working elsewhere. We want to stop somewhere that makes sense.
// Was thinking that we can give the user the control of what are the paths that we
// should stop for.
func PullConfigFiles(stopDirs map[string]bool) (*map[string]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting the current working directory: %s", err)
		return &map[string]string{}, err
	}

	configPathsChannel := make(chan string)
	go func() {
		currentDir := cwd
		for {
			present, configPath := ContainsConfigFile(currentDir)
			fmt.Println(present, configPath)
			if present {
				configPathsChannel <- configPath
			}

			if stopDirs[currentDir] {
				fmt.Printf("stopDirs contains currentDir: %s\n", currentDir)
				break
			}

			currentDir = filepath.Dir(currentDir)
			fmt.Printf("Moving up the directory tree: %s\n", currentDir)
		}

		close(configPathsChannel)
	}()

	for path := range configPathsChannel {
		fmt.Printf("path from channel: %s\n", path)
	}

	return &map[string]string{}, nil
}

func IsConfigFile(fileName string) bool {
	return fileName == "dyncomp.json"
}

func ContainsConfigFile(directory string) (bool, string) {
	configPath := filepath.Join(directory, "dyncomp.json")
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
	json.Unmarshal(contents, &configFileMap)
	return configFileMap, nil
}
