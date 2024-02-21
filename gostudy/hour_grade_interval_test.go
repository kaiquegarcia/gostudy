package gostudy_test

import (
	"testing"
	"time"

	"github.com/kaiquegarcia/gostudy/gostudy"
	"github.com/stretchr/testify/assert"
)

func Test_HourGradeInterval(t *testing.T) {
	start, _ := time.Parse(gostudy.LayoutTimeOnly, "10:00")
	end, _ := time.Parse(gostudy.LayoutTimeOnly, "12:00")

	t.Run("should extend end from 12:00 to 15:00", func(t *testing.T) {
		// Arrange
		hgi := &gostudy.HourGradeInterval{
			Start: start,
			End:   end,
		}

		startOnEnd, _ := time.Parse(gostudy.LayoutTimeOnly, "12:00")
		endToExtend, _ := time.Parse(gostudy.LayoutTimeOnly, "15:00")

		// Act
		result := hgi.Extends(startOnEnd, endToExtend)

		// Assert
		assert.True(t, result, "extends should return true")
		assert.Equal(t, start, hgi.Start, "start should not change")
		assert.Equal(t, endToExtend, hgi.End, "end should change to endToExtend")
	})

	t.Run("should extend start from 10:00 to 08:00", func(t *testing.T) {
		// Arrange
		hgi := &gostudy.HourGradeInterval{
			Start: start,
			End:   end,
		}

		startToExtend, _ := time.Parse(gostudy.LayoutTimeOnly, "08:00")
		endOnStart, _ := time.Parse(gostudy.LayoutTimeOnly, "10:00")

		// Act
		result := hgi.Extends(startToExtend, endOnStart)

		// Assert
		assert.True(t, result, "extends should return true")
		assert.Equal(t, startToExtend, hgi.Start, "start should change to startToExtend")
		assert.Equal(t, end, hgi.End, "end should not change")
	})

	t.Run("should override start and end", func(t *testing.T) {
		// Arrange
		hgi := &gostudy.HourGradeInterval{
			Start: start,
			End:   end,
		}

		startToExtend, _ := time.Parse(gostudy.LayoutTimeOnly, "08:00")
		endToExtend, _ := time.Parse(gostudy.LayoutTimeOnly, "15:00")

		// Act
		result := hgi.Extends(startToExtend, endToExtend)

		// Assert
		assert.True(t, result, "extends should return true")
		assert.Equal(t, startToExtend, hgi.Start, "start should change to startToExtend")
		assert.Equal(t, endToExtend, hgi.End, "end should change to endToExtend")
	})

	t.Run("should not extend", func(t *testing.T) {
		// Arrange
		hgi := &gostudy.HourGradeInterval{
			Start: start,
			End:   end,
		}

		startAfterEnd, _ := time.Parse(gostudy.LayoutTimeOnly, "15:00")
		otherEnd, _ := time.Parse(gostudy.LayoutTimeOnly, "16:00")

		// Act
		result := hgi.Extends(startAfterEnd, otherEnd)

		// Assert
		assert.False(t, result, "extends should return false")
		assert.Equal(t, start, hgi.Start, "start should not change")
		assert.Equal(t, end, hgi.End, "end should not change")
	})
}
