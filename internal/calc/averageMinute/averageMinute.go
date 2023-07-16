package averageminute

// Process and functions responsible for taking in events and calculating the average duration per minute
// These averages will then be sent to the Moving Average process

import (
	"event_cli/internal/data"
)

// workers are thread objects that receive events and calculate the average duration for them.
type worker struct {
	commandChannel <-chan command             // a channel for receiving commands from the master thread
	resultChannel  chan<- data.RunningAverage // a channel for periodically outputting the worker's results
}

// commands received by workers, emitted by the master thread
type command struct {
	cmdType       commandType // The type of the command
	eventArgument *data.Event // Only used for the ProcessEvent command
}

// Enumerate of all the types of command
type commandType int

const (
	ProcessEvent commandType = iota // command that sends an event to be processed by a worker thread
	EmitResults                     // command that asks all worker threads to return what they have gathered so far
	Cancel                          // command to stop workers in case of fatal error
)

// Handler function for each worker. Reveives events and calculates their average
//
//	Once master thread asks for results back, the worker thread outputs their running average count and resets it
func (w worker) startWorker() {

	var runningAverage data.RunningAverage

	for command := range w.commandChannel {
		switch command.cmdType {
		case ProcessEvent:
			{
				runningAverage.AddMeasurment(command.eventArgument.Duration)
			}
		case EmitResults:
			{
				w.resultChannel <- runningAverage
				runningAverage.Reset()
			}
		case Cancel:
			{
				return
			}
		}
	}

}

// Main process for the Average per minute calculation. This process spawns worker threads in order to parallelize
// the average calculation.
// Parallelization is achieved by spreading events between worker threads and having each thread calculate their own local average.
// Once an event is received that has a timestamp with a different minute, master thread knows that all events for the current minute
// have been received so it asks all worker threads to emit their results. Then it performs a weighted average of all
// the running averages returned by the workers, obtaining the average of the last minute. It then bundles the average
// with an identifier of that minute and sends it to the next stage.
func StartAveragePerMinuteProcess(numWorkers int, eventChannel <-chan *data.Event) <-chan data.MinuteAverage {

	var currMinute data.UnixMinute

	//channel that will receive each worker's resulting averages
	resultChannel := make(chan data.RunningAverage, numWorkers*2) //channel buffers scale linearly with number of worker threads

	//object responsible for spreading (fanning out) commands between all workers
	fanOutManager := newFanOut[command](numWorkers, 2)

	//channel that will output final averages for each minute. This channel will be listened to by the next stage
	minuteChannel := make(chan data.MinuteAverage) // channel used to send all the averages per minute

	// Create the worker threads
	for i := 0; i < numWorkers; i++ {
		w := worker{
			commandChannel: fanOutManager.outputChannels[i],
			resultChannel:  resultChannel,
		}
		go w.startWorker()
	}

	// Create the master thread
	go func() {
		defer close(minuteChannel)
		defer fanOutManager.close()
		defer close(resultChannel)
		for event := range eventChannel {

			// New minute found. Ask threads to show their work and reset their running averages
			if currMinute != data.GetUnixMinuteFromTime(event.Timestamp.Time) {

				// Edge case: If initial minute is uninitialized send an empty average, this will be the first line of the file
				if currMinute == 0 {
					currMinute = data.GetUnixMinuteFromTime(event.Timestamp.Time)

					minuteChannel <- data.MinuteAverage{
						Minute:  currMinute,
						Average: data.RunningAverage{},
					}

				} else {

					// broadcast an emit results command
					fanOutManager.broadcast(command{
						cmdType: EmitResults,
					})

					// gather worker's results
					minuteChannel <- gatherResults(resultChannel, numWorkers, currMinute)

					// update current minute to match most recent event
					currMinute = data.GetUnixMinuteFromTime(event.Timestamp.Time)
				}

			}

			//send out the newest event to a random thread
			fanOutManager.send(command{
				cmdType:       ProcessEvent,
				eventArgument: event,
			})
		}

		//No more events, get the last average
		fanOutManager.broadcast(command{
			cmdType: EmitResults,
		})

		minuteChannel <- gatherResults(resultChannel, numWorkers, currMinute)
	}()

	return minuteChannel

}

// this function receives the results from all threads, performs a weighed average between all threads and packages it into a data.MinuteAverage object
func gatherResults(resultChannel <-chan data.RunningAverage, numWorkers int, currentMinute data.UnixMinute) data.MinuteAverage {
	minuteAverage := data.MinuteAverage{
		Minute:  currentMinute + 1,
		Average: data.RunningAverage{},
	}
	var resultingRunningAverage data.RunningAverage
	for i := 0; i < numWorkers; i++ {
		resultingRunningAverage = <-resultChannel
		minuteAverage.Average.AddRunningAverage(resultingRunningAverage)
	}

	return minuteAverage

}
