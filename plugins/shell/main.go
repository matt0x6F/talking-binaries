package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mattouille/talking-binaries/protocol"
)

// readFromPipe reads from a pipe
func readFromPipe() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}

		output = append(output, input)
	}

	return []byte(string(output)), nil
}

type Config struct {
	Commands []string
	Parallel bool
}

type Input struct {
	Config Config
}

func main() {
	logBuff := bytes.NewBuffer(nil)
	logger := log.New(logBuff, "[plugin-shell] ", log.Flags())

	logger.Printf("Beginning execution")

	data, err := readFromPipe()
	if err != nil {
		logger.Printf("Error while reading from pipe: %s", err)

		os.Exit(1)
	}

	if data == nil {
		logger.Printf("No input received")

		os.Exit(1)
	}

	var input Input

	err = json.Unmarshal(data, &input)
	if err != nil {
		logger.Printf("Error unmarshaling input: %s", err)

		os.Exit(1)
	}

	logger.Printf("Config received: %+v", input)

	output := protocol.NewOutput()

	logger.Printf("Preparing to execute %d commands", len(input.Config.Commands))

	for i, command := range input.Config.Commands {
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)

		parts := strings.Split(command, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		logger.Printf("Executing command: %s", cmd.String())

		if err := cmd.Run(); err != nil {
			logger.Printf("Error encountered while executing command: %s", err)
		}

		output.UserOutput[i] = stdout.Bytes()
		output.Errors[i] = stderr.Bytes()
	}

	output.PluginLogs = logBuff.Bytes()

	data, err = json.Marshal(output)
	if err != nil {
		logger.Printf("Error while marshaling output: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", data)
}
