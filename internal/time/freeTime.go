package time

import (
	"fmt"
	"time"
)

type FreeTime struct {
	hours   int
	minutes int
}

func NewFreeTime(hours, minutes int) FreeTime {
	return FreeTime{
		hours,
		minutes,
	}
}

func NewFreeTimeFromDate(date time.Time) FreeTime {
	return FreeTime{
		hours:   date.Hour(),
		minutes: date.Minute(),
	}
}

func (f *FreeTime) ToString() string {
	return fmt.Sprintf("%02d:%02d", f.hours, f.minutes)
}

func (f *FreeTime) GetMinutes() int {
	return f.hours*60 + f.minutes
}

func (f *FreeTime) Add(min int) {
	minutes := f.minutes + min
	if minutes < 60 {
		f.minutes = minutes
		return
	}
	hours := minutes / 60
	if minutes%60 != 0 || minutes == 60 {
		f.minutes = minutes - hours*60
	}
	f.hours += hours
	return
}

func (f *FreeTime) Elapsed(end *FreeTime, stepMin int) (steps int) {
	min := f.GetMinutes() - end.GetMinutes()
	steps = min / stepMin
	return
}

func (f *FreeTime) IsLess(time *FreeTime) bool {
	return f.GetMinutes() < time.GetMinutes()
}

func GetFreeTime(star, end, now FreeTime, stepMin int) []FreeTime {
	steps := end.Elapsed(&star, stepMin)
	freeTimes := make([]FreeTime, 0, steps+1)
	if !star.IsLess(&now) {
		freeTimes = append(freeTimes, star)
	}

	for i := 0; i < steps; i++ {
		star.Add(stepMin)
		if star.IsLess(&now) {
			continue
		}
		freeTimes = append(freeTimes, star)
	}

	return freeTimes
}
