package timeutil

import "time"

const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day
)

var Now = time.Now
