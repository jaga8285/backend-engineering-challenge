package ioprocess

import (
	"bufio"
	"event_cli/internal/data"
	"fmt"
	"os"
)

const eventBufferSize = 20

func OpenFile(filename string) (<-chan *data.Event, <-chan error, error) {

	outputChannel := make(chan *data.Event, eventBufferSize)
	errorChannel := make(chan error, 1)

	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)
	go func() {
		defer close(outputChannel)
		defer close(errorChannel)
		defer file.Close()
		for scanner.Scan() {
			event, err := data.EventFromString(scanner.Text())
			if err != nil {
				fmt.Print("heey")
				errorChannel <- err
				continue
			}
			outputChannel <- event
		}
		if scanner.Err() != nil {
			fmt.Print("hey")
			errorChannel <- scanner.Err()
		}
	}()
	return outputChannel, errorChannel, nil
}
