package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mattouille/talking-binaries/config"
	"github.com/mattouille/talking-binaries/logger"
	"github.com/mattouille/talking-binaries/protocol"
)

var (
	plugins = map[string]string{}
)

func main() {
	log := logger.New(os.Stdout, true)

	cfg, err := config.LoadFromDisk()
	if err != nil {
		log.Info("Error while loading config: %s", err)

		os.Exit(1)
	}

	log.Debug("Got config: %+v", cfg)

	// validate plugin configuration against plugins found on the system
	for i, plugin := range cfg.Plugins {
		if plugin.Name == "" {
			log.Info("Plugin location cannot be empty")

			os.Exit(1)
		}

		if plugin.Path == "" {
			// will fail if the binary is not executable
			path, err := exec.LookPath(plugin.Name)
			if err != nil {
				log.Info("Unable to determine path for plugin %s [%s]: %s", plugin.Ref, plugin.Name, err)

				os.Exit(1)
			}

			log.Debug("Found plugin at location: %s", path)

			abs, err := filepath.Abs(path)
			if err != nil {
				log.Info("Unable to determine absolute path for plugin: %s", err)

				os.Exit(1)
			}

			plugin.Path = abs
			cfg.Plugins[i].Path = abs
		} else {
			_, err := os.Stat(plugin.Path)
			if err != nil {
				log.Info("Error while validating plugin location: %s", err)
			}
		}

		plugins[plugin.Ref] = plugin.Path
	}

	log.Debug("Built plugin store: %+v", plugins)

	timeTrack := func(start time.Time, name string) {
		elapsed := time.Since(start)
		log.Debug("%s took %s", name, elapsed)
	}

	// step execution
	for _, step := range cfg.Steps {
		fullExecStart := time.Now()

		log.Info("Executing step: %s", step.Name)

		rawInput := protocol.Input{
			Config:   step.Config,
			Parallel: false,
		}

		jsonInput, err := json.Marshal(rawInput)
		if err != nil {
			log.Info("Error preparing plugin configuration: %s", err)

			os.Exit(1)
		}

		log.Debug("Sending data: %s", string(jsonInput))

		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)

		cmd := exec.Command(plugins[step.Plugin])
		cmd.Stdin = bytes.NewBuffer(jsonInput)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		binaryExecutionTime := time.Now()

		err = cmd.Run()
		if err != nil {
			log.Info("Error encountered while executing command: %s", err)
		}

		timeTrack(binaryExecutionTime, "Binary execution time")

		var output protocol.Output

		err = json.Unmarshal(stdout.Bytes(), &output)
		if err != nil {
			log.Info("Error encountered while unmarshaling command output: %s", err)
		}

		fmt.Fprintf(os.Stdout, string(output.PluginLogs))

		fmt.Fprintf(os.Stdout, "User Output:\n")

		if output.UserOutput == nil {
			fmt.Fprintf(os.Stdout, "None\n")
		}

		for i := 0; i < len(output.UserOutput); i++ {
			fmt.Fprintf(os.Stdout, string(output.UserOutput[i]))

			fmt.Fprintf(os.Stdout, "\n")
		}

		timeTrack(fullExecStart, "Step execution time")
	}
}
