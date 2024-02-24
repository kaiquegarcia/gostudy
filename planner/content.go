package planner

import (
	"time"
)

type Content struct {
	Subject   string
	Title     string
	Duration  time.Duration
	Reference string
	Attempts  int
}

func newContentFromRow(columns []string) (*Content, error) {
	// Subject, Title, Duration, Reference
	if len(columns) != 4 {
		return nil, ErrUnexpectedColumnsLength
	}

	duration, err := parseDuration(columns[2])
	if err != nil {
		return nil, err
	}

	return &Content{
		Subject:   columns[0],
		Title:     columns[1],
		Duration:  duration,
		Reference: columns[3],
		Attempts:  0,
	}, nil
}

func (c *Content) IsBetween(start time.Time, end time.Time) bool {
	return !start.Add(c.Duration).After(end)
}
