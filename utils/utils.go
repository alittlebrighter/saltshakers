package utils

import (
	"strconv"
	"time"
)

// SecondOfTheMonth returns midnight of the second day of the month @arg date is in.
// The second day is chosen to avoid timezone issues.
func SecondOfTheMonth(date time.Time) time.Time {
	yearStr := strconv.Itoa(date.Year())
	monthStr := strconv.Itoa(int(date.Month()))
	secondOfTheMonth, _ := time.Parse("1/2/2006", monthStr+"/2/"+yearStr)
	return secondOfTheMonth
}
