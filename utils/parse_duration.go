package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidDurationFormat = fmt.Errorf("the duration doesn't follow the hh:mm:ss pattern")

func ParseDuration(duration string) (time.Duration, error) {
	pieces := strings.Split(duration, ":")
	if len(pieces) != 3 {
		return 0, ErrInvalidDurationFormat
	}
	for _, p := range pieces {
		if len(p) != 2 {
			return 0, ErrInvalidDurationFormat
		}
	}

	hours, err := strconv.Atoi(pieces[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(pieces[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(pieces[2])
	if err != nil {
		return 0, err
	}

	return (time.Duration(hours) * time.Hour) + (time.Duration(minutes) * time.Minute) + (time.Duration(seconds) * time.Second), nil
}
