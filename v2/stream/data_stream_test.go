package stream_test

import (
	"os"
	"testing"

	"github.com/kaiquegarcia/gostudy/v2/stream"
	"github.com/stretchr/testify/assert"
)

func Test_CSVDataStream(t *testing.T) {
	file, _ := os.Open("file_for_test.csv")
	defer file.Close()

	t.Run("should return headers on first call", func(t *testing.T) {
		// Arrange
		cds, err := stream.NewCSVDataStream(file)
		if !assert.Nil(t, err, "err from NewCSVDataStream should be nil") {
			t.FailNow()
		}

		// Act
		row, err := cds.Read()
		if !assert.Nil(t, err, "err from Read should be nil") {
			t.FailNow()
		}

		// Assert
		if !assert.Len(t, row, 4, "row should have 4 columns") {
			t.FailNow()
		}
		assert.Equal(t, "Header 1", row[0], "row[0] should be 'Header 1'")
		assert.Equal(t, "Header 2", row[1], "row[1] should be 'Header 2'")
		assert.Equal(t, "Header 3", row[2], "row[2] should be 'Header 3'")
		assert.Equal(t, "Header 4", row[3], "row[3] should be 'Header 4'")
	})

	t.Run("should return row 1 on second call", func(t *testing.T) {
		// Arrange
		cds, err := stream.NewCSVDataStream(file)
		if !assert.Nil(t, err, "err from NewCSVDataStream should be nil") {
			t.FailNow()
		}

		// Act
		_, err = cds.Read()
		if !assert.Nil(t, err, "err from first Read should be nil") {
			t.FailNow()
		}

		row, err := cds.Read()
		if !assert.Nil(t, err, "err from second Read should be nil") {
			t.FailNow()
		}

		// Assert
		if !assert.Len(t, row, 4, "row should have 4 columns") {
			t.FailNow()
		}
		assert.Equal(t, "A1", row[0], "row[0] should be 'A1'")
		assert.Equal(t, "B1", row[1], "row[1] should be 'B1'")
		assert.Equal(t, "C1", row[2], "row[2] should be 'C1'")
		assert.Equal(t, "D1", row[3], "row[3] should be 'D1'")
	})

	t.Run("should return row 1 on first call after Unread second call", func(t *testing.T) {
		// Arrange
		cds, err := stream.NewCSVDataStream(file)
		if !assert.Nil(t, err, "err from NewCSVDataStream should be nil") {
			t.FailNow()
		}

		// Act
		_, err = cds.Read()
		if !assert.Nil(t, err, "err from first Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from second Read should be nil") {
			t.FailNow()
		}

		err = cds.Unread()
		if !assert.Nil(t, err, "err from Unread should be nil") {
			t.FailNow()
		}

		row, err := cds.Read()
		if !assert.Nil(t, err, "err from Read-after-Unread should be nil") {
			t.FailNow()
		}

		// Assert
		if !assert.Len(t, row, 4, "row should have 4 columns") {
			t.FailNow()
		}
		assert.Equal(t, "A1", row[0], "row[0] should be 'A1'")
		assert.Equal(t, "B1", row[1], "row[1] should be 'B1'")
		assert.Equal(t, "C1", row[2], "row[2] should be 'C1'")
		assert.Equal(t, "D1", row[3], "row[3] should be 'D1'")
	})

	t.Run("should return row 3 after 4 reads", func(t *testing.T) {
		// Arrange
		cds, err := stream.NewCSVDataStream(file)
		if !assert.Nil(t, err, "err from NewCSVDataStream should be nil") {
			t.FailNow()
		}

		// Act
		_, err = cds.Read()
		if !assert.Nil(t, err, "err from first Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from second Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from third Read should be nil") {
			t.FailNow()
		}

		row, err := cds.Read()
		if !assert.Nil(t, err, "err from fourth Read should be nil") {
			t.FailNow()
		}

		// Assert
		if !assert.Len(t, row, 4, "row should have 4 columns") {
			t.FailNow()
		}
		assert.Equal(t, "A3", row[0], "row[0] should be 'A3'")
		assert.Equal(t, "B3", row[1], "row[1] should be 'B3'")
		assert.Equal(t, "C3", row[2], "row[2] should be 'C3'")
		assert.Equal(t, "D3", row[3], "row[3] should be 'D3'")
	})

	t.Run("should return EOF after 5 reads", func(t *testing.T) {
		// Arrange
		cds, err := stream.NewCSVDataStream(file)
		if !assert.Nil(t, err, "err from NewCSVDataStream should be nil") {
			t.FailNow()
		}

		// Act
		_, err = cds.Read()
		if !assert.Nil(t, err, "err from first Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from second Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from third Read should be nil") {
			t.FailNow()
		}

		_, err = cds.Read()
		if !assert.Nil(t, err, "err from fourth Read should be nil") {
			t.FailNow()
		}

		row, err := cds.Read()

		// Assert
		assert.ErrorIs(t, err, stream.ErrEOF, "err should be EOF")
		assert.Len(t, row, 0, "row should have 0 columns")
	})
}
