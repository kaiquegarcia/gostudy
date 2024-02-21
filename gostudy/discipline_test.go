package gostudy_test

import (
	"testing"
	"time"

	"github.com/kaiquegarcia/gostudy/gostudy"
	"github.com/stretchr/testify/assert"
)

func Test_Discipline(t *testing.T) {
	cs1 := []*gostudy.Content{
		{Title: "S1-C1", Duration: 10 * time.Minute},
		{Title: "S1-C2", Duration: 10 * time.Minute},
	}
	cs2 := []*gostudy.Content{
		{Title: "S2-C1", Duration: 10 * time.Minute},
		{Title: "S2-C2", Duration: 10 * time.Minute},
	}

	t.Run("first next should return first subject and first content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, "S1", subject.Name, "subject name should be S1")
		assert.Equal(t, "S1-C1", content.Title, "content title should be S1-C1")
	})

	t.Run("second next should return first subject and second content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err from second Next() should be nil")
		assert.Equal(t, "S1", subject.Name, "subject name should be S1")
		assert.Equal(t, "S1-C2", content.Title, "content title should be S1-C2")
	})

	t.Run("third next should return second subject and first content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err from third Next() should be nil")
		assert.Equal(t, "S2", subject.Name, "subject name should be S2")
		assert.Equal(t, "S2-C1", content.Title, "content title should be S2-C1")
	})

	t.Run("fourth next should return second subject and second content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from third Next() should be nil")
			t.FailNow()
		}

		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err from fourth Next() should be nil")
		assert.Equal(t, "S2", subject.Name, "subject name should be S2")
		assert.Equal(t, "S2-C2", content.Title, "content title should be S2-C2")
	})

	t.Run("fifth next should return ErrEndOfList", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from third Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from fourth Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()

		// Assert
		assert.ErrorIs(t, gostudy.ErrEndOfList, err, "err from fifth Next() should be ErrEndOfList")
	})

	t.Run("next after back after first next should return first subject and first content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		err = discipline.Back()
		if err != nil {
			assert.Nil(t, err, "err from Back() should be nil")
			t.FailNow()
		}

		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err from second Next() should be nil")
		assert.Equal(t, "S1", subject.Name, "subject name should be S1")
		assert.Equal(t, "S1-C1", content.Title, "content title should be S1-C1")
	})

	t.Run("next after back after third next should return second subject and first content", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		_, _, err := discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from first Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from second Next() should be nil")
			t.FailNow()
		}

		_, _, err = discipline.Next()
		if err != nil {
			assert.Nil(t, err, "err from third Next() should be nil")
			t.FailNow()
		}

		err = discipline.Back()
		if err != nil {
			assert.Nil(t, err, "err from Back() should be nil")
			t.FailNow()
		}

		content, subject, err := discipline.Next()

		// Assert
		assert.Nil(t, err, "err from fourth Next() should be nil")
		assert.Equal(t, "S2", subject.Name, "subject name should be S2")
		assert.Equal(t, "S2-C1", content.Title, "content title should be S2-C1")
	})

	t.Run("first back should return ErrBeginningOfList", func(t *testing.T) {
		// Arrange
		ss1 := gostudy.NewSubject("S1")
		ss1.SetContents(cs1)
		ss2 := gostudy.NewSubject("S2")
		ss2.SetContents(cs2)
		discipline := gostudy.NewDiscipline("Math", "math.csv", 0, 0, 0)
		discipline.SetSubjects([]*gostudy.Subject{ss1, ss2})

		// Act
		err := discipline.Back()

		// Assert
		assert.ErrorIs(t, gostudy.ErrBeginningOfList, err, "err should be ErrBeginningOfList")
	})
}
