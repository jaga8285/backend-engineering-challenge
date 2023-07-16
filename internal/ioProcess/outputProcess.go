package ioprocess

import (
	"encoding/json"
	"event_cli/internal/data"
	"fmt"
	"io"
	"os"
)

func WriteFile(filename string, averageChannel <-chan data.CalculatedMovingAverage) (chan error, error) {

	var f io.Writer
	if filename == "" {
		f = os.Stdout
	} else {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		f = file
	}

	errorChannel := make(chan error)

	go func() {

		defer close(errorChannel)

		if file, ok := f.(*os.File); ok {
			defer file.Close()
		}

		for average := range averageChannel {
			line, err := json.Marshal(average)
			if err != nil {
				errorChannel <- err
			}
			fmt.Fprintln(f, string(line))
		}

	}()

	return errorChannel, nil
}
