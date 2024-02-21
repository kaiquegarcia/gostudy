package gostudy_test

import (
	"testing"
	"time"

	"github.com/kaiquegarcia/gostudy/gostudy"
	"github.com/stretchr/testify/assert"
)

func Test_Subject(t *testing.T) {
	contents := []*gostudy.Content{
		{Title: "C1", Duration: 10 * time.Minute},
		{Title: "C2", Duration: 10 * time.Minute},
		{Title: "C3", Duration: 10 * time.Minute},
		{Title: "C4", Duration: 10 * time.Minute},
		{Title: "C5", Duration: 10 * time.Minute},
	}

	t.Run("next should retrieve first content", func(t *testing.T) {
		// Arrange
		subject := gostudy.NewSubject("")
		subject.SetContents(contents)

		// Act
		content, err := subject.Next()
		if err != nil {
			assert.Nil(t, err, "err should be nil")
			t.FailNow()
		}

		// Assert
		assert.Equal(t, "C1", content.Title)
	})

	t.Run("next should retrieve second content", func(t *testing.T) {
		// Arrange
		subject := gostudy.NewSubject("")
		subject.SetContents(contents)

		// Act
		// first call
		_, err := subject.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		// second call
		content, err := subject.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		// Assert
		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, "C2", content.Title)
	})

	t.Run("next after back should retrieve first content", func(t *testing.T) {
		// Arrange
		subject := gostudy.NewSubject("")
		subject.SetContents(contents)

		// Act
		// first call
		_, err := subject.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		// back call
		err = subject.Back()
		if err != nil {
			assert.Nil(t, err, "err from Back() should be nil")
			t.FailNow()
		}

		// second call
		content, err := subject.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		// Assert
		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, "C1", content.Title)
	})

	t.Run("back before first next should return error", func(t *testing.T) {
		// Arrange
		subject := gostudy.NewSubject("")
		subject.SetContents(contents)

		// Act
		// first call
		err := subject.Back()

		// Assert
		assert.ErrorIs(t, gostudy.ErrBeginningOfList, err, "err should be ErrBeginningOfList")
	})

	t.Run("next after last number of contents should return error", func(t *testing.T) {
		// Arrange
		subject := gostudy.NewSubject("")
		subject.SetContents(contents)

		// Act
		for i := 0; i < len(contents); i++ {
			c, err := subject.Next()
			if err != nil {
				assert.Nil(t, err, "err from call '%d' Next() should be nil", i)
				t.FailNow()
			}

			assert.Equal(t, contents[i].Title, c.Title, "title should be equal")
		}

		_, err := subject.Next()

		// Assert
		assert.ErrorIs(t, gostudy.ErrEndOfList, err, "err should be ErrEndOfList")
	})
}
