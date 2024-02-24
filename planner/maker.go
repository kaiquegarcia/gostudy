package planner

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/kaiquegarcia/gostudy/v2/logging"
	"github.com/kaiquegarcia/gostudy/v2/stream"
)

type Maker struct {
	outputFile                   *os.File
	outputWriter                 *csv.Writer
	hg                           HourGrade
	disciplines                  []*Discipline
	inputedStartDate             time.Time
	logger                       logging.Logger
	checkedDisciplinesCount      int
	currentDisciplineIndex       int
	currentDayDisciplineDuration time.Duration
	finishedDisciplinesIndexes   []int
}

func NewMaker(
	logger logging.Logger,
	hg HourGrade,
	data []*Discipline,
	startDate time.Time,
	outputFilename string,
) (*Maker, error) {
	file, err := os.Create(outputFilename)
	if err != nil {
		return nil, err
	}

	cw := csv.NewWriter(file)
	cw.Write([]string{"Datetime", "Discipline", "Subject", "Title", "Reference", "Duration"})
	cw.Flush()
	return &Maker{
		hg:                         hg,
		disciplines:                data,
		inputedStartDate:           startDate,
		logger:                     logger,
		currentDisciplineIndex:     -1,
		checkedDisciplinesCount:    0,
		finishedDisciplinesIndexes: make([]int, 0),
		outputFile:                 file,
		outputWriter:               cw,
	}, nil
}

func (p *Maker) Close() {
	p.outputWriter.Flush()
	p.outputFile.Close()
	for _, d := range p.disciplines {
		d.Close()
	}
}

func (p *Maker) hasExploredAllDisciplines() bool {
	return p.checkedDisciplinesCount >= len(p.disciplines)
}

func (p *Maker) isAllDisciplinesFinished() bool {
	return len(p.finishedDisciplinesIndexes) == len(p.disciplines)
}

func (p *Maker) isDisciplineFinished(disciplineIndex int) bool {
	for _, index := range p.finishedDisciplinesIndexes {
		if index == disciplineIndex {
			return true
		}
	}

	return false
}

func (p *Maker) nextDiscipline(hgi *HourGradeInterval) {
	if p.currentDayDisciplineDuration > 0 {
		gap := p.disciplines[p.currentDisciplineIndex].SubjectGap
		p.logger.Debug("discipline has duration > 0, adding subject gap of %s", gap)
		hgi.Start = hgi.Start.Add(gap)
	}

	p.currentDisciplineIndex++
	p.currentDisciplineIndex %= len(p.disciplines)
	p.currentDayDisciplineDuration = 0
	p.checkedDisciplinesCount++
}

func (p *Maker) startDate() (time.Time, error) {
	if !p.hg.HasGradeFor(p.inputedStartDate) {
		p.logger.Debug("start date doesn't have hour on the grade, getting next date")
		date, err := p.hg.NextDate(p.inputedStartDate)
		if err != nil {
			p.logger.Error(err, "could not retrieve next date")
			return time.Time{}, err
		}

		p.logger.Debug("next date retrieved successfully: %s", date.Format(LayoutDateOnly))
		return date, nil
	}

	p.logger.Debug(
		"start date has hour on the grade, using it as current date: %s",
		p.inputedStartDate.Format(LayoutDateOnly),
	)
	return p.inputedStartDate, nil
}

func (p *Maker) Mount() error {
	date, err := p.startDate()
	if err != nil {
		return err
	}

	p.currentDisciplineIndex = 0
	p.logger.Debug("starting mount loop")
	for {
		err = p.mountDate(date)
		if err == stream.ErrEOF {
			break
		}

		if err != nil {
			return err
		}

		p.logger.Debug("intervals loop finished, getting next date")
		date, err = p.hg.NextDate(date)
		if err != nil {
			p.logger.Error(err, "could not retrieve next date")
			return err
		}

		p.logger.Debug("next date retrieved successfully: %s, moving to the loop", date.Format(LayoutDateOnly))
	}

	return nil
}
