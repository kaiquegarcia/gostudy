package gostudy

import (
	"time"

	"github.com/kaiquegarcia/gostudy/utils"
)

type Planner struct {
	hg                           HourGrade
	disciplines                  []*Discipline
	inputedStartDate             time.Time
	logger                       utils.Logger
	checkedDisciplinesCount      int
	currentDisciplineIndex       int
	currentDayDisciplineDuration time.Duration
	finishedDisciplinesIndexes   []int
	output                       []PlannerOutput
}

func NewPlanner(
	logger utils.Logger,
	hg HourGrade,
	data []*Discipline,
	startDate time.Time,
) *Planner {
	return &Planner{
		hg:                         hg,
		disciplines:                data,
		inputedStartDate:           startDate,
		logger:                     logger,
		currentDisciplineIndex:     -1,
		checkedDisciplinesCount:    0,
		finishedDisciplinesIndexes: make([]int, 0),
		output:                     make([]PlannerOutput, 0),
	}
}

func (p *Planner) hasExploredAllDisciplines() bool {
	return p.checkedDisciplinesCount >= len(p.disciplines)
}

func (p *Planner) isAllDisciplinesFinished() bool {
	return len(p.finishedDisciplinesIndexes) == len(p.disciplines)
}

func (p *Planner) isDisciplineFinished(disciplineIndex int) bool {
	for _, index := range p.finishedDisciplinesIndexes {
		if index == disciplineIndex {
			return true
		}
	}

	return false
}

func (p *Planner) nextDiscipline(hgi *HourGradeInterval) {
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

func (p *Planner) startDate() (time.Time, error) {
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

func (p *Planner) Mount() error {
	date, err := p.startDate()
	if err != nil {
		return err
	}

	p.output = make([]PlannerOutput, 0)
	p.currentDisciplineIndex = 0
	p.logger.Debug("starting mount loop")
	for {
		err = p.mountDate(date)
		if err == ErrEndOfList {
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

func (p *Planner) ResultRecords() [][]string {
	records := make([][]string, len(p.output)+1)
	records[0] = []string{"Datetime", "Discipline", "Subject", "Title", "Reference", "Duration"}
	for index, plannerOutput := range p.output {
		records[index+1] = plannerOutput.ToRecord()
	}

	return records
}
