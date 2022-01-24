package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/mattouille/talking-binaries/plugin"
	"github.com/mattouille/talking-binaries/step"
)

const (
	DefaultRepoConfigName = "DevFile"
)

type RepoConfig struct {
	// Plugins that can be used throughout the steps.
	Plugins []plugin.Object
	// An ordered execution of steps.
	Steps []step.Object
}

// LoadFromDisk will search for configuration in the current directory. If no file is found then it searches recursively
// lower until it meets root.
func LoadFromDisk() (RepoConfig, error) {
	var (
		config RepoConfig
	)

	// get the working directory
	dir, err := os.Getwd()
	if err != nil {
		return RepoConfig{}, fmt.Errorf("unable to determine working directory: %w", err)
	}

	path, err := searchForConfig(dir, DefaultRepoConfigName)
	if err != nil {
		return RepoConfig{}, err
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return RepoConfig{}, err
	}

	// read in config
	if err = yaml.Unmarshal(contents, &config); err != nil {
		log.Printf("Error reading YAML config: %s", err)

		return RepoConfig{}, err
	}

	return config, nil
}

// searchForConfig traverses down a directory structure looking for a config file. If none is found then no error is
// returned along with a blank path.
func searchForConfig(dir, name string) (string, error) {
	// get the files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("unable to read directory: %w", err)
	}

	// search for the config file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Name() == name {
			return filepath.Abs(file.Name())
		}
	}

	path := filepath.Dir(dir)

	// usually true at the root directory; Dir of "/" is "/"
	if path == dir {
		return "", fmt.Errorf("config file not found")
	}

	return searchForConfig(path, name)
}
