package movingaverage

import (
	"event_cli/internal/data"
)

// Main Process for moving average calculation. This process will receive minute averages from the average per minute process and output
// the average for the past n minutes, where n is window size. This process assumes that minute averages are ascending in timestamp
// but it does not assume that all minutes will be provided. Namely, it is capable of "filling in" missing minutes with empty values.
// The process stores the past n minute averages in a circular buffer and sends the weighted average along with a timestamp to the file output process
func StartMovingAverage(windowSize int, averagesChannel <-chan data.MinuteAverage) <-chan data.CalculatedMovingAverage {
	movingAverageChannel := make(chan data.CalculatedMovingAverage)
	lastNAverages := data.NewFIFOQueue[data.MinuteAverage](windowSize)
	var lastUnixMinute data.UnixMinute

	go func() {
		defer close(movingAverageChannel)
		average := <-averagesChannel

		lastUnixMinute = average.Minute
		lastNAverages.Enqueue(average)
		movingAverageChannel <- data.CalcMovingAverage(lastNAverages.GetQueue(), average.Minute)

		for average := range averagesChannel {
			for lastUnixMinute < average.Minute-1 {
				lastUnixMinute++
				lastNAverages.Enqueue(data.MinuteAverage{
					Minute: lastUnixMinute,
				})
				avg := data.CalcMovingAverage(lastNAverages.GetQueue(), lastUnixMinute)

				movingAverageChannel <- avg

			}

			lastUnixMinute = average.Minute
			lastNAverages.Enqueue(average)
			avg := data.CalcMovingAverage(lastNAverages.GetQueue(), average.Minute)
			movingAverageChannel <- avg
		}

	}()

	return movingAverageChannel
}
