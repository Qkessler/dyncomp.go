package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: Think about where to stop, since it shouldn't work with the HOME directory
// when working elsewhere. We want to stop somewhere that makes sense.
// Was thinking that we can give the user the control of what are the paths that we
// should stop for.
func PullConfigFiles(stopDirs map[string]bool) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error while getting the current working directory: %s", err)
		return []string{}, err
	}

	// FIXME: This is currently walking downwards.
	var configPaths []string
	filepath.WalkDir(cwd, func(path string, dirEntry fs.DirEntry, err error) error {
		fmt.Println(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		fileInfo, err := dirEntry.Info()
		if err != nil {
			fmt.Println(err)
			return err
		} else if fileInfo.IsDir() && path != cwd {
			return fs.SkipDir
		}

		if stopDirs[path] {
			configPath := filepath.Join(path, "dyncomp.json")
			if _, err := os.Stat(configPath); err == nil {
				configPaths = append(configPaths, configPath)
			}

			return fs.SkipDir
		}
		if IsConfigFile(fileInfo.Name()) {
			configPaths = append(configPaths, path)
			return nil
		}

		return err
	})
	return configPaths, nil
}

func IsConfigFile(fileName string) bool {
	return fileName == "dyncomp.json"
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
