package planner

import (
	"os"
	"time"

	"github.com/kaiquegarcia/gostudy/v2/stream"
)

type Discipline struct {
	Name          string
	Filename      string
	DailyLimit    time.Duration
	ContentGap    time.Duration
	SubjectGap    time.Duration
	contentStream stream.DataStream
}

func NewDiscipline(
	name string,
	filename string,
	dailyLimit time.Duration,
	contentGap time.Duration,
	subjectGap time.Duration,
) (*Discipline, error) {
	contentFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	contentStream, err := stream.NewCSVDataStream(contentFile)
	if err != nil {
		return nil, err
	}

	// sets pointer to first row, skipping header
	_, err = contentStream.Read()
	if err != nil {
		return nil, err
	}

	return &Discipline{
		Name:          name,
		Filename:      filename,
		DailyLimit:    dailyLimit,
		ContentGap:    contentGap,
		SubjectGap:    subjectGap,
		contentStream: contentStream,
	}, nil
}

func (d *Discipline) Close() error {
	return d.contentStream.Close()
}

func (d *Discipline) Next() (*Content, error) {
	columns, err := d.contentStream.Read()
	if err != nil {
		return nil, err
	}

	return newContentFromRow(columns)
}

func (d *Discipline) Back() error {
	return d.contentStream.Unread()
}

func NewDisciplineFromRows(rows [][]string) ([]*Discipline, error) {
	disciplines := make([]*Discipline, 0)
	for line := 1; line < len(rows); line++ {
		columns := rows[line]
		if len(columns) == 0 {
			break
		}
		// Name, Filename, Daily Limit, Content Gap, Subject Gap
		if len(columns) != 5 {
			return nil, ErrUnexpectedColumnsLength
		}

		dailyLimit, err := parseDuration(columns[2])
		if err != nil {
			return nil, err
		}

		contentGap, err := parseDuration(columns[3])
		if err != nil {
			return nil, err
		}

		subjectGap, err := parseDuration(columns[4])
		if err != nil {
			return nil, err
		}

		discipline, err := NewDiscipline(
			columns[0],
			columns[1],
			dailyLimit,
			contentGap,
			subjectGap,
		)
		if err != nil {
			return nil, err
		}

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}
