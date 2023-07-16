package averageminute

// Fan out object! Capable of receiving an object T and sending it to one of multiple channels or broadcasting it to all channels.
// (This is a very rudimentary implementation)
type fanOut[T any] struct {
	outputChannels    []chan T
	lastChannelSentTo int
}

// Send message. Sends the message to a random, non-full channel
// WARNING: very bad implementation, don't look at it
func (fan fanOut[T]) send(message T) {
	var sent bool
	numChannels := len(fan.outputChannels)
	for !sent {
		fan.lastChannelSentTo = (fan.lastChannelSentTo + 1) % numChannels
		select {
		case fan.outputChannels[fan.lastChannelSentTo] <- message:
			sent = true
		default:
		}
	}
}

func (fan fanOut[T]) broadcast(message T) {
	for i := 0; i < len(fan.outputChannels); i++ {
		fan.outputChannels[i] <- message
	}
}
func (fan fanOut[T]) close() {
	for i := 0; i < len(fan.outputChannels); i++ {
		close(fan.outputChannels[i])
	}
}

func newFanOut[T any](numChanels int, bufferSize int) fanOut[T] {
	chs := make([]chan T, numChanels)

	for i := 0; i < numChanels; i++ {
		ch := make(chan T, bufferSize)
		chs[i] = ch
	}
	return fanOut[T]{
		outputChannels: chs,
	}
}
