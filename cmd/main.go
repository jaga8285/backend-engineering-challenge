package main

import (
	averageminute "event_cli/internal/calc/averageMinute"
	movingaverage "event_cli/internal/calc/movingAverage"
	ioprocess "event_cli/internal/ioProcess"
	"fmt"
	"os"

	"github.com/fred1268/go-clap/clap"
)

type config struct {
	InputFile  string `clap:"--input_file,-i,mandatory"` // Input file where the events are stored
	OutputFile string `clap:"--output_file,-o"`          // Output file where averages will be stored. If left empty, write to stdout
	WindowSize int    `clap:"--window_size,mandatory"`   // The size of the window of the moving average
	NumWorkers int    `clap:"--num_threads"`             // Number of thread used during the average calculation step (default:4)
}

func main() {

	cfg := &config{
		NumWorkers: 4,
	}
	if _, err := clap.Parse(os.Args, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nUsage:event_cli --window-size N --input_file file [--output_file file --num_threads N]\n", err)
		return
	}
	run(*cfg)
}

// The algorithm is split into 4 stages, all of the stages run concurrently and share information between channels

func run(cfg config) {

	//First stage is File input. Receives a file and returns a channel that outputs events

	eventChannel, _, err := ioprocess.OpenFile(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file:%v\n", err)
		return
	}

	//  Second stage is average by minute calculation. This stage receives events and outputs the average of the last minute received,
	//once all events for that minute have been received
	//  MinuteAverage is a struct that functions as a bucket storing the average event duration for a given minute in time.
	minuteAverageChannel := averageminute.StartAveragePerMinuteProcess(cfg.NumWorkers, eventChannel)

	//  Third stage is moving average caluclation. This stage receives the average by minute returned by the previous stage
	//and calculates the average of the past N minutes, where N is the window size
	//  This stage outputs moving averages and their corresponding timestamp (CalculatedMovingAverage)
	calculatedAveragesChannel := movingaverage.StartMovingAverage(cfg.WindowSize, minuteAverageChannel)

	// The Fourth and final stage receives the CalculatedMovingAverages output by the previous stage and writes them to the destination file or stdout
	writeFileErrorChannel, err := ioprocess.WriteFile(cfg.OutputFile, calculatedAveragesChannel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file:%v\n", err)
		return
	}

	for err := range writeFileErrorChannel {
		fmt.Fprintf(os.Stderr, "Error during file write: %v", err)

	}
}
