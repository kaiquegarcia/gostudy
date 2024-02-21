package gostudy

import (
	"sort"
	"strings"
	"time"
)

const (
	LayoutTimeOnly         = "15:04"
	LayoutTimeWithTimezone = "15:04:05Z07:00"
	LayoutDateOnly         = "2006-01-02"
)

type HourGrade map[time.Weekday][]*HourGradeInterval

func NewHourGrade() HourGrade {
	hg := HourGrade{}
	for weekday := 0; weekday < 7; weekday++ {
		hg[time.Weekday(weekday)] = make([]*HourGradeInterval, 0)
	}

	return hg
}

func (hg HourGrade) Add(weekday time.Weekday, start time.Time, end time.Time) {
	var extendedInterval *HourGradeInterval = nil
	for index, hgi := range hg[weekday] {
		if hgi.Extends(start, end) {
			hg[weekday] = append(hg[weekday][:index], hg[weekday][index+1:]...)
			extendedInterval = hgi
			break
		}
	}

	if extendedInterval != nil {
		hg.Add(weekday, extendedInterval.Start, extendedInterval.End)
		return
	}

	hg[weekday] = append(hg[weekday], &HourGradeInterval{Start: start, End: end})
}

func (hg HourGrade) Sort() {
	for _, intervals := range hg {
		if len(intervals) == 0 {
			continue
		}

		sort.Slice(intervals, func(i, j int) bool {
			return !intervals[i].Start.After(intervals[j].End)
		})
	}
}

func (hg HourGrade) HasGradeFor(date time.Time) bool {
	return len(hg[date.Weekday()]) > 0
}

func (hg HourGrade) NextDate(from time.Time) (time.Time, error) {
	// clone to be able to check if we come back to the same weekday
	date := from.Add(0)
	for {
		date = date.AddDate(0, 0, 1)
		if hg.HasGradeFor(date) {
			return date, nil
		}

		if date.Weekday() == from.Weekday() {
			break
		}
	}

	return time.Time{}, ErrUnavailableWeekdays
}

func (hg HourGrade) IntervalsFor(date time.Time) ([]*HourGradeInterval, error) {
	intervals := make([]*HourGradeInterval, len(hg[date.Weekday()]))
	for index, hgi := range hg[date.Weekday()] {
		start, err := hgi.SetStartTime(date)
		if err != nil {
			return nil, err
		}

		end, err := hgi.SetEndTime(date)
		if err != nil {
			return nil, err
		}

		intervals[index] = &HourGradeInterval{
			Start: start,
			End:   end,
		}
	}

	return intervals, nil
}

func ExtractHourGradeFromTableRecords(records [][]string) (HourGrade, error) {
	if len(records) < 8 {
		return nil, ErrUnexpectedGradeLength
	}

	hg := NewHourGrade()
	for weekdayOffset := 1; weekdayOffset < 8; weekdayOffset++ {
		columns := records[weekdayOffset]
		for columnIndex := 1; columnIndex < len(columns); columnIndex++ {
			entry := columns[columnIndex]
			if entry == "" {
				break
			}

			entryData := strings.Split(entry, "-")
			if len(entryData) != 2 {
				return nil, ErrUnexpectedIntervalLength
			}

			start, err := time.Parse(LayoutTimeOnly, entryData[0])
			if err != nil {
				return nil, err
			}

			end, err := time.Parse(LayoutTimeOnly, entryData[1])
			if err != nil {
				return nil, err
			}

			hg.Add(time.Weekday(weekdayOffset-1), start, end)
		}
	}

	hg.Sort()
	return hg, nil
}
