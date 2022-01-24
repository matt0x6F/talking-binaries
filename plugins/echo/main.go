package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	data, err := readFromPipe()
	if err != nil {
		log.Printf("Error while reading from pipe: %s", err)

		os.Exit(1)
	}

	if data == nil {
		log.Printf("No data received")

		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Configuration received: %s\n", data)
}
