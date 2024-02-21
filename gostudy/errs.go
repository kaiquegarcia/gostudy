package gostudy

import "fmt"

var (
	ErrBeginningOfList           = fmt.Errorf("it's on the beginning of the list already")
	ErrEndOfList                 = fmt.Errorf("reached end of the list")
	ErrUnexpectedIntervalLength  = fmt.Errorf("the time interval must have only two elements, the beginning and the end of the interval")
	ErrUnexpectedGradeLength     = fmt.Errorf("the hour grade spreadsheet must have at least 7 rows, one row per day of week")
	ErrUnexpectedColumnsLength   = fmt.Errorf("the columns number doesn't match with the required count")
	ErrUnavailableWeekdays       = fmt.Errorf("could not find any weekday with available hour grade intervals")
	ErrContentDurationUnplayable = fmt.Errorf("content duration is unplayable")
)
