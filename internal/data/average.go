package data

import (
	"encoding/json"
	"math/big"
	"time"
)

type RunningAverage struct {
	sum     uint
	counter uint
}

func (rm RunningAverage) GetAverage() float64 {
	if rm.sum == 0 && rm.counter == 0 {
		return 0
	}
	return float64(rm.sum) / float64(rm.counter)
}

func (rm *RunningAverage) AddMeasurment(measurment uint) {
	rm.sum += measurment
	rm.counter++
}

func (rm *RunningAverage) AddRunningAverage(otherAverage RunningAverage) {
	rm.sum = rm.sum + otherAverage.sum
	rm.counter = rm.counter + otherAverage.counter
}

func (rm *RunningAverage) Reset() {
	rm.sum = 0
	rm.counter = 0
}

// UnixTime counted in minutes: a time represented by the number of minutes elapsed since January 1, 1970 UTC
type UnixMinute int

func GetUnixMinuteFromTime(time time.Time) UnixMinute {
	return UnixMinute(time.Unix() / 60)
}
func GetTimeFromUnixMinute(minute UnixMinute) time.Time {
	return time.Unix(int64(minute)*60, 0)
}

// MinuteAverage stores the average of values of minute interval, and identifies that minute
type MinuteAverage struct {
	Minute  UnixMinute
	Average RunningAverage
}

type CalculatedMovingAverage struct {
	Minute  UnixMinute
	Average float64
}

const AverageTimeLayout = "2006-01-02 15:04:05"

func (average CalculatedMovingAverage) MarshalJSON() ([]byte, error) {

	dateString := GetTimeFromUnixMinute(average.Minute).Format(AverageTimeLayout)

	return json.Marshal(&struct {
		Date                  string  `json:"date"`
		Average_delivery_time float64 `json:"average_delivery_time"`
	}{
		Date:                  dateString,
		Average_delivery_time: average.Average,
	})
}

func CalcMovingAverage(averages []MinuteAverage, currentMinute UnixMinute) CalculatedMovingAverage {
	sum := big.NewInt(0)
	counter := big.NewInt(0)
	for _, average := range averages {
		if average.Average.counter == 0 {
			continue
		}
		averageSum := new(big.Int).SetUint64(uint64(average.Average.sum))
		averageCount := new(big.Int).SetUint64(uint64(average.Average.counter))
		sum.Add(sum, averageSum)
		counter.Add(counter, averageCount)
	}

	if sum.Cmp(big.NewInt(0)) == 0 && counter.Cmp(big.NewInt(0)) == 0 {
		return CalculatedMovingAverage{
			Minute:  currentMinute,
			Average: 0,
		}
	}

	sumFloat := new(big.Float).SetInt(sum)
	counterFloat := new(big.Float).SetInt(counter)

	movingAverage, _ := new(big.Float).Quo(sumFloat, counterFloat).Float64()
	return CalculatedMovingAverage{
		Minute:  currentMinute,
		Average: movingAverage,
	}
}
