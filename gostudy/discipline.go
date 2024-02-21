package gostudy

import (
	"time"

	"github.com/kaiquegarcia/gostudy/utils"
)

type Discipline struct {
	Name            string
	Filename        string
	DailyLimit      time.Duration
	ContentGap      time.Duration
	SubjectGap      time.Duration
	subjects        []*Subject
	currentPointer  int
	previousPointer int
}

func NewDiscipline(
	name string,
	filename string,
	dailyLimit time.Duration,
	contentGap time.Duration,
	subjectGap time.Duration,
) *Discipline {
	return &Discipline{
		Name:            name,
		Filename:        filename,
		DailyLimit:      dailyLimit,
		ContentGap:      contentGap,
		SubjectGap:      subjectGap,
		subjects:        make([]*Subject, 0),
		currentPointer:  0,
		previousPointer: 0,
	}
}

func (d *Discipline) Add(s *Subject) {
	d.subjects = append(d.subjects, s)
}

func (d *Discipline) SetSubjects(ss []*Subject) {
	d.subjects = ss
	d.currentPointer = 0
	d.previousPointer = 0
}

func (d *Discipline) Next() (*Content, *Subject, error) {
	// try to get next from current subject
	s := d.subjects[d.currentPointer]
	c, err := s.Next()
	if err == nil {
		d.previousPointer = d.currentPointer
		return c, s, nil
	}

	if err != ErrEndOfList {
		return nil, nil, err
	}

	if d.currentPointer >= len(d.subjects)-1 {
		return nil, nil, ErrEndOfList
	}

	d.previousPointer = d.currentPointer
	d.currentPointer++
	s = d.subjects[d.currentPointer]
	c, err = s.Next()
	if err == nil {
		return c, s, nil
	}

	return nil, nil, err
}

func (d *Discipline) Back() error {
	s := d.subjects[d.currentPointer]
	err := s.Back()
	if err != nil {
		return err
	}

	if d.previousPointer < d.currentPointer {
		d.currentPointer--
		d.previousPointer--
	}

	return nil
}

func (d *Discipline) TotalDuration() time.Duration {
	t := time.Duration(0)
	for _, s := range d.subjects {
		t += s.TotalDuration()
	}

	return t
}

func ExtractDisciplineContentFromTableRecords(
	discipline *Discipline,
	records [][]string,
) error {
	subjectIndexes := make(map[string]int)
	subjects := make([]*Subject, 0)
	for index := 1; index < len(records); index++ {
		columns := records[index] // [Subject, Title, Duration, Reference]
		if len(columns) == 0 {
			// end of table
			break
		}

		if len(columns) < 4 {
			return ErrUnexpectedColumnsLength
		}

		duration, err := utils.ParseDuration(columns[2])
		if err != nil {
			return err
		}

		var subject *Subject = nil
		if subjectIndex, exists := subjectIndexes[columns[0]]; exists {
			subject = subjects[subjectIndex]
		} else {
			subject = NewSubject(columns[0])
			subjectIndexes[columns[0]] = len(subjects)
			subjects = append(subjects, subject)
		}

		subject.Add(&Content{
			Title:     columns[1],
			Duration:  duration,
			Reference: columns[3],
			Attempts:  0,
		})
	}

	discipline.SetSubjects(subjects)
	return nil
}

func ExtractDisciplineListFromTableRecords(records [][]string) ([]*Discipline, error) {
	disciplines := make([]*Discipline, len(records)-1)
	for index := 1; index < len(records); index++ {
		columns := records[index] // Name, Filename, Daily Limit, Content Gap, Subject Gap
		if len(columns) == 0 {
			// end of table
			break
		}

		if len(columns) < 5 {
			return nil, ErrUnexpectedColumnsLength
		}

		dailyLimit, err := utils.ParseDuration(columns[2])
		if err != nil {
			return nil, err
		}

		contentGap, err := utils.ParseDuration(columns[3])
		if err != nil {
			return nil, err
		}

		subjectGap, err := utils.ParseDuration(columns[4])
		if err != nil {
			return nil, err
		}

		disciplines[index-1] = NewDiscipline(
			columns[0],
			columns[1],
			dailyLimit,
			contentGap,
			subjectGap,
		)
	}

	return disciplines, nil
}
