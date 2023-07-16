package movingaverage_test

import (
	movingaverage "event_cli/internal/calc/movingAverage"
	"event_cli/internal/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRunningAverage(sum uint, counter uint) data.RunningAverage {
	var average data.RunningAverage
	average.AddMeasurment(sum)
	for i := uint(1); i < counter; i++ {
		average.AddMeasurment(0)
	}
	return average
}

var minuteAverageExample1 = []data.MinuteAverage{{
	Minute:  11,
	Average: *new(data.RunningAverage),
}, {
	Minute:  12,
	Average: newRunningAverage(20, 1),
}, {
	Minute:  16,
	Average: newRunningAverage(31, 1),
}, {
	Minute:  24,
	Average: newRunningAverage(54, 1),
},
}

var minuteAverageExample2 = []data.MinuteAverage{{
	Minute:  0,
	Average: *new(data.RunningAverage),
}, {
	Minute:  1,
	Average: newRunningAverage(10, 2),
}, {
	Minute:  2,
	Average: newRunningAverage(15, 4),
}, {
	Minute:  3,
	Average: newRunningAverage(7, 3),
}, {
	Minute:  4,
	Average: newRunningAverage(1, 1),
}, {
	Minute:  5,
	Average: newRunningAverage(15, 8),
}, {
	Minute:  6,
	Average: newRunningAverage(30, 1),
}, {
	Minute:  7,
	Average: newRunningAverage(17, 5),
}, {
	Minute:  8,
	Average: newRunningAverage(11, 6),
},
}

func TestAverages1(t *testing.T) {

	targetMinuteAverages := []float64{0, 20, 20, 20, 20, 25.5, 25.5, 25.5, 25.5, 25.5, 25.5, 31, 31, 42.5}

	windowSize := 10

	minuteAverageChannel := make(chan data.MinuteAverage)

	calculatedAveragesChannel := movingaverage.StartMovingAverage(windowSize, minuteAverageChannel)

	go func() {
		for _, e := range minuteAverageExample1 {
			minuteAverageChannel <- e
		}
		close(minuteAverageChannel)
	}()

	i := 0
	for calculatedAverage := range calculatedAveragesChannel {
		assert.Equal(t, targetMinuteAverages[i], calculatedAverage.Average, "Target %v has a mismatch: %v", i, calculatedAverage.Average)
		i++
	}
}

func TestAverages2(t *testing.T) {

	targetMinuteAverages := []float64{0, 5, 4.166666666666667, 3.5555555555555554, 2.875, 1.9166666666666667, 4.6, 4.428571428571429, 4.833333333333333}

	windowSize := 3

	minuteAverageChannel := make(chan data.MinuteAverage)

	calculatedAveragesChannel := movingaverage.StartMovingAverage(windowSize, minuteAverageChannel)

	go func() {
		for _, e := range minuteAverageExample2 {
			minuteAverageChannel <- e
		}
		close(minuteAverageChannel)
	}()

	i := 0
	for calculatedAverage := range calculatedAveragesChannel {
		assert.Equal(t, targetMinuteAverages[i], calculatedAverage.Average, "Target %v has a mismatch: %v", i, calculatedAverage.Average)
		i++
	}
}

/*
func TestMinutes(t *testing.T) {

	targetMinutes := []data.UnixMinute{1689367812 / 60, 1689367872 / 60, 1689367992 / 60, 1689368052 / 60, 1689368112 / 60}

	eventChannel := make(chan *data.Event)

	minuteAverageChannel := averageminute.StartAveragePerMinuteProcess(2, eventChannel)

	go func() {
		for _, e := range eventsExample {
			eventChannel <- e
		}
		close(eventChannel)
	}()

	i := 0
	for minuteAverage := range minuteAverageChannel {
		assert.Equal(t, targetMinutes[i], minuteAverage.Average.GetAverage(), "Target %v has a mismatch: %v", i, minuteAverage.Average)
		i++
	}
}*/
