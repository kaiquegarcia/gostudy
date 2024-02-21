package gostudy

import (
	"fmt"
	"time"
)

type HourGradeInterval struct {
	Start time.Time
	End   time.Time
}

func (hgi *HourGradeInterval) Extends(start time.Time, end time.Time) bool {
	extended := false
	if !start.After(hgi.Start) && !end.Before(hgi.Start) {
		hgi.Start = start
		extended = true
	}

	if !start.After(hgi.End) && !end.Before(hgi.End) {
		hgi.End = end
		extended = true
	}

	return extended
}

func (hgi *HourGradeInterval) SetStartTime(date time.Time) (time.Time, error) {
	return time.Parse(
		time.RFC3339,
		fmt.Sprintf(
			"%sT%s",
			date.Format(LayoutDateOnly),
			hgi.Start.Format(LayoutTimeWithTimezone),
		),
	)
}

func (hgi *HourGradeInterval) SetEndTime(date time.Time) (time.Time, error) {
	return time.Parse(
		time.RFC3339,
		fmt.Sprintf(
			"%sT%s",
			date.Format(LayoutDateOnly),
			hgi.End.Format(LayoutTimeWithTimezone),
		),
	)
}
