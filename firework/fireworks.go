package firework

import (
	"time"
)

// CronFirework :
type CronFirework struct {
	*DelayFirework
	Delta    time.Duration
	Seconds  []int
	Minutes  []int
	Hours    []int
	Days     []int
	Months   []int
	Years    []int
	IndexArr []int
}

// Copy :
func (f *CronFirework) Copy() IFirework {
	return &CronFirework{
		Seconds:       f.Seconds,
		Minutes:       f.Minutes,
		Hours:         f.Hours,
		Days:          f.Days,
		Months:        f.Months,
		Years:         f.Years,
		DelayFirework: f.DelayFirework.Copy().(*DelayFirework),
	}
}

func (f *CronFirework) getTimeValues() []int {
	return []int{
		f.Time.Second(),
		f.Time.Minute(),
		f.Time.Hour(),
		f.Time.Day(),
		int(f.Time.Month()),
		f.Time.Year(),
	}
}

func (f *CronFirework) getTimeRanges() [][]int {
	return [][]int{
		f.Seconds,
		f.Minutes,
		f.Hours,
		f.Days,
		f.Months,
		f.Years,
	}
}

// UpdateIndex :
func (f *CronFirework) UpdateIndex() {
	timeValues := f.getTimeValues()
	timeRanges := f.getTimeRanges()
	f.IndexArr = make([]int, len(timeValues))
	for i, v := range timeValues {
		values := timeRanges[i]
		for ii, vv := range values {
			if v == vv {
				f.IndexArr[i] = ii
			}
		}
	}
}

// UpdateTime :
func (f *CronFirework) UpdateTime() bool {
	timeValues := f.getTimeValues()
	timeRanges := f.getTimeRanges()
	base := 1
	for i, index := range f.IndexArr {
		values := timeRanges[i]
		index += base
		valuesLen := len(values)
		if index >= valuesLen {
			base = 1
			index = index % valuesLen
		} else {
			base = 0
		}
		f.IndexArr[i] = index
		timeValues[i] = values[index]
	}
	f.Time = time.Date(
		timeValues[5],
		time.Month(timeValues[4]),
		timeValues[3],
		timeValues[2],
		timeValues[1],
		timeValues[0],
		0, time.Local,
	)
	for i, v := range f.getTimeValues() {
		if timeValues[i] != v {
			return f.UpdateTime()
		}
	}
	return true
}
