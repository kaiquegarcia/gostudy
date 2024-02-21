package gostudy

import "time"

type Content struct {
	Title     string
	Duration  time.Duration
	Reference string
}

func (c *Content) IsBetween(start time.Time, end time.Time) bool {
	return !start.Add(c.Duration).After(end)
}
