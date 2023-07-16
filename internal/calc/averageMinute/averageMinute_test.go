package averageminute_test

import (
	averageminute "event_cli/internal/calc/averageMinute"
	"event_cli/internal/data"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const numWorkers = 4

var eventsExample1 = []*data.Event{{
	Timestamp: data.UnmarshalerTime{Time: time.Date(2023, 07, 14, 20, 50, 12, 0, time.UTC)},
	Duration:  1,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2023, 07, 14, 20, 51, 12, 0, time.UTC)},
	Duration:  2,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2023, 07, 14, 20, 51, 20, 0, time.UTC)},
	Duration:  3,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2023, 07, 14, 20, 53, 12, 0, time.UTC)},
	Duration:  4,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2023, 07, 14, 20, 54, 12, 0, time.UTC)},
	Duration:  5,
}}

var eventsExample2 = []*data.Event{{
	Timestamp: data.UnmarshalerTime{Time: time.Date(2018, 12, 26, 18, 11, 8, 0, time.UTC)},
	Duration:  20,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2018, 12, 26, 18, 15, 19, 0, time.UTC)},
	Duration:  31,
}, {
	Timestamp: data.UnmarshalerTime{Time: time.Date(2018, 12, 26, 18, 23, 19, 0, time.UTC)},
	Duration:  54,
}}

func TestAverages1(t *testing.T) {

	targetAverages := []float64{0, 1, 2.5, 4, 5}

	auxTestAverage(t, eventsExample1, targetAverages)
}
func TestMinutes1(t *testing.T) {

	targetMinutes := []data.UnixMinute{1689367812 / 60, 1689367872 / 60, 1689367932 / 60, 1689368052 / 60, 1689368112 / 60}

	auxTestMinutes(t, eventsExample1, targetMinutes)

}

func TestAverages2(t *testing.T) {

	targetAverages := []float64{0, 20, 31, 54}

	auxTestAverage(t, eventsExample2, targetAverages)
}
func TestMinutes2(t *testing.T) {

	targetMinutes := []data.UnixMinute{1545847868 / 60, 1545847928 / 60, 1545848179 / 60, 1545848659 / 60}

	auxTestMinutes(t, eventsExample2, targetMinutes)

}

func auxTestAverage(t *testing.T, inputEvents []*data.Event, target []float64) {

	eventChannel := make(chan *data.Event)

	minuteAverageChannel := averageminute.StartAveragePerMinuteProcess(numWorkers, eventChannel)

	go func() {
		for _, e := range inputEvents {
			eventChannel <- e
		}
		close(eventChannel)
	}()

	i := 0
	for minuteAverage := range minuteAverageChannel {
		assert.Equal(t, target[i], minuteAverage.Average.GetAverage(), "Target %v has a mismatch: %+v", i, minuteAverage.Average)
		i++
	}
}

func auxTestMinutes(t *testing.T, inputEvents []*data.Event, target []data.UnixMinute) {
	eventChannel := make(chan *data.Event)

	minuteAverageChannel := averageminute.StartAveragePerMinuteProcess(numWorkers, eventChannel)

	go func() {
		for _, e := range inputEvents {
			eventChannel <- e
		}
		close(eventChannel)
	}()

	i := 0
	for minuteAverage := range minuteAverageChannel {
		assert.Equal(t, target[i], minuteAverage.Minute, "Target %v has a mismatch", i)
		i++
	}
}
