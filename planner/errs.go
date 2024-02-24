package planner

import "fmt"

var (
	ErrInvalidDurationFormat     = fmt.Errorf("the duration doesn't follow the hh:mm:ss pattern")
	ErrUnavailableWeekdays       = fmt.Errorf("could not find any weekday with available hour grade intervals")
	ErrUnexpectedColumnsLength   = fmt.Errorf("the columns number doesn't match with the required count")
	ErrUnexpectedIntervalLength  = fmt.Errorf("the time interval must have only two elements, the beginning and the end of the interval")
	ErrUnexpectedGradeLength     = fmt.Errorf("the hour grade spreadsheet must have at least 7 rows, one row per day of week")
	ErrContentDurationUnplayable = fmt.Errorf("content duration is unplayable")
)
