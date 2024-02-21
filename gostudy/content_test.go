package gostudy_test

import (
	"testing"
	"time"

	"github.com/kaiquegarcia/gostudy/gostudy"
	"github.com/stretchr/testify/assert"
)

func Test_Content(t *testing.T) {
	content := &gostudy.Content{
		Duration: 30 * time.Minute,
	}

	t.Run("IsBetween should match with margin", func(t *testing.T) {
		// Arrange
		start, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
		end, _ := time.Parse(time.RFC3339, "2024-01-01T11:00:00Z") // 11:00 - 10:00 = 1 hour

		// Act
		result := content.IsBetween(start, end)

		// Assert
		assert.True(t, result, "result should be true")
	})

	t.Run("IsBetween should match perfectly", func(t *testing.T) {
		// Arrange
		start, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
		end, _ := time.Parse(time.RFC3339, "2024-01-01T10:30:00Z") // 10:30 - 10:00 = 30 minutes

		// Act
		result := content.IsBetween(start, end)

		// Assert
		assert.True(t, result, "result should be true")
	})

	t.Run("IsBetween should not match even for 1 second", func(t *testing.T) {
		// Arrange
		start, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
		end, _ := time.Parse(time.RFC3339, "2024-01-01T10:29:59Z") // 10:29:59 - 10:00 = 29 minutes 59 seconds

		// Act
		result := content.IsBetween(start, end)

		// Assert
		assert.False(t, result, "result should be false")
	})
}
