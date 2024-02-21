package gostudy

import (
	"time"

	"github.com/kaiquegarcia/gostudy/utils"
)

type Planner struct {
	HourGrade                    HourGrade
	Disciplines                  []*Discipline
	StartDate                    time.Time
	checkedDisciplinesCount      int
	currentDisciplineIndex       int
	currentDayDisciplineDuration time.Duration
	finishedDisciplinesIndexes   []int
	output                       []PlannerOutput
}

func NewPlanner(
	hg HourGrade,
	data []*Discipline,
	start time.Time,
) *Planner {
	return &Planner{
		HourGrade:                  hg,
		Disciplines:                data,
		StartDate:                  start,
		currentDisciplineIndex:     -1,
		checkedDisciplinesCount:    0,
		finishedDisciplinesIndexes: make([]int, 0),
		output:                     make([]PlannerOutput, 0),
	}
}

func (p *Planner) hasExploredAllDisciplines() bool {
	return p.checkedDisciplinesCount >= len(p.Disciplines)
}

func (p *Planner) isAllDisciplinesFinished() bool {
	return len(p.finishedDisciplinesIndexes) == len(p.Disciplines)
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
		hgi.Start = hgi.Start.Add(p.Disciplines[p.currentDisciplineIndex].SubjectGap)
	}

	p.currentDisciplineIndex++
	p.currentDisciplineIndex %= len(p.Disciplines)
	p.currentDayDisciplineDuration = 0
	p.checkedDisciplinesCount++
}

func (p *Planner) Mount(logger utils.Logger) error {
	date := p.StartDate
	if !p.HourGrade.HasGradeFor(date) {
		logger.Debug("start date doesn't have hour on the grade, getting next date")
		var err error
		date, err = p.HourGrade.NextDate(p.StartDate)
		if err != nil {
			logger.Error(err, "could not retrieve next date")
			return err
		}

		logger.Debug("next date retrieved successfully: %s", date.Format(LayoutDateOnly))
	} else {
		logger.Debug("start date has hour on the grade, using it as current date: %s", date.Format(LayoutDateOnly))
	}

	p.output = make([]PlannerOutput, 0)
	p.currentDisciplineIndex = 0
	logger.Debug("starting mount loop")
	for {
		dateStr := date.Format(LayoutDateOnly)
		logger.Debug("retrieving time intervals for date %s", dateStr)
		intervals, err := p.HourGrade.IntervalsFor(date)
		if err != nil {
			logger.Error(err, "could not retrieve time intervals for date %s", dateStr)
			return err
		}

		p.currentDayDisciplineDuration = 0
		logger.Debug("%d intervals found, start loop", len(intervals))
		for _, hgi := range intervals {
			initialStr := hgi.Start.Format(LayoutTimeOnly)
			endStr := hgi.End.Format(LayoutTimeOnly)
			logger.Debug("starting procedure for interval %s-%s", initialStr, endStr)
			var (
				previousDiscipline *Discipline
				previousSubject    *Subject
				isFirst            = true
			)
			p.checkedDisciplinesCount = 0

			loopCounter := 0
			for {
				loopCounter++
				startStr := hgi.Start.Format(LayoutTimeOnly)
				logger.Debug("inner loop %d, checking if interval is still able to proceed", loopCounter)
				if !hgi.Start.Before(hgi.End) {
					logger.Debug("already reached the end of the interval, breaking at inner loop %d", loopCounter)
					p.checkedDisciplinesCount = 0
					break
				}

				logger.Debug(
					"inner loop %d still have time left (%s - %s), checking if already explored all possibilities",
					loopCounter, startStr, endStr,
				)
				if p.hasExploredAllDisciplines() {
					logger.Debug("already explored all possibilities, breaking at inner loop %d", loopCounter)
					p.checkedDisciplinesCount = 0
					break
				}

				logger.Debug("still have disciplines to explore, checking current discipline index %d", p.currentDisciplineIndex)
				discipline := p.Disciplines[p.currentDisciplineIndex]
				logger.Debug("which means '%s' ~ checking if it's already finished", discipline.Name)
				if p.isDisciplineFinished(p.currentDisciplineIndex) {
					logger.Debug("discipline '%s' is already finished, adding gap only if current duration is higher than zero", discipline.Name)
					p.nextDiscipline(hgi)
					previousDiscipline = discipline
					continue
				}

				logger.Debug("discipline '%s' it's not finished, checking if it already exhausted daily limit (%s)", discipline.Name, discipline.DailyLimit)
				if p.currentDayDisciplineDuration >= discipline.DailyLimit {
					logger.Debug("discipline '%s' already exhausted daily limit, adding gap (if duration is higher than zero)", discipline.Name)
					p.nextDiscipline(hgi)
					previousDiscipline = discipline
					continue
				}

				logger.Debug("discipline '%s' did not exhaust daily limit, checking next content", discipline.Name)
				content, subject, err := discipline.Next()
				if err == ErrEndOfList {
					logger.Debug("discipline '%s' reached the end of content list, checking if all disciplines finished", discipline.Name)
					p.finishedDisciplinesIndexes = append(p.finishedDisciplinesIndexes, p.currentDisciplineIndex)

					if p.isAllDisciplinesFinished() {
						logger.Debug("all disciplines finished, ending planner mount! last discipline = %s", discipline.Name)
						return nil
					}

					logger.Debug("still have disciplines to work on. getting next discipline, adding gap only if current duration is higher than zero")
					p.nextDiscipline(hgi)
					continue
				}

				if err != nil {
					logger.Error(err, "could not get next discipline content")
					return err
				}

				logger.Debug("discipline's content retrieved. checking if we should include a gap before the content")
				totalDuration := content.Duration
				var preGap time.Duration = 0
				if isFirst {
					logger.Debug("it's the first content of this time interval, no gap is required")
					preGap = 0
				} else if subject != previousSubject && previousSubject != nil && previousDiscipline == discipline {
					totalDuration += discipline.SubjectGap
					preGap = discipline.SubjectGap
					logger.Debug("it's a new subject of the same discipline, adding gap of %s", preGap)
				} else if subject == previousSubject {
					totalDuration += discipline.ContentGap
					preGap = discipline.ContentGap
					logger.Debug("it's a new content of the same subject, adding gap of %s", preGap)
				} else {
					logger.Debug("it's a new content from other disciplines, no gap is required")
				}

				logger.Debug("checking if discipline current duration (%s) + totalDuration (%s) exhaust discipline daily limit (%s)", p.currentDayDisciplineDuration, totalDuration, discipline.DailyLimit)
				if (p.currentDayDisciplineDuration + totalDuration) > discipline.DailyLimit {
					logger.Debug("discipline '%s's gap exhausts daily limit, getting next discipline, adding gap (if duration is higher than zero)", discipline.Name)
					p.nextDiscipline(hgi)
					previousDiscipline = discipline
					continue
				}

				logger.Debug("discipline '%s's gap doesn't exhaust daily limit", discipline.Name)

				logger.Debug("checking if the content + gap (%s) can be added to time interval (%s-%s)", totalDuration, startStr, endStr)
				if hgi.End.Sub(hgi.Start) < totalDuration {
					logger.Debug("too much duration for the time left, calling discipline.Back()")
					err = discipline.Back()
					if err != nil {
						logger.Error(err, "could not step back on the discipline's content")
						return err
					}

					logger.Debug("as discipline '%s' can't fill with the current content, we'll call the next discipline", discipline.Name)
					p.nextDiscipline(hgi)
					continue
				}

				logger.Debug("content can be added to time interval, adding to PlannerOutput list")
				p.output = append(p.output, PlannerOutput{
					Time:       hgi.Start.Add(preGap),
					Discipline: discipline,
					Subject:    subject,
					Content:    content,
				})

				previousSubject = subject
				previousDiscipline = discipline
				p.currentDayDisciplineDuration += totalDuration
				hgi.Start = hgi.Start.Add(totalDuration)
				isFirst = false
				logger.Debug("inner loop %d finished, starting next", loopCounter)
			}

			logger.Debug("procedure for interval  %s-%s finished, calling next interval", initialStr, endStr)
		}

		logger.Debug("intervals loop finished, getting next date")
		date, err = p.HourGrade.NextDate(date)
		if err != nil {
			logger.Error(err, "could not retrieve next date")
			return err
		}

		logger.Debug("next date retrieved successfully: %s, moving to the loop", date.Format(LayoutDateOnly))
	}
}

func (p *Planner) ResultRecords() [][]string {
	records := make([][]string, len(p.output)+1)
	records[0] = []string{"Datetime", "Discipline", "Subject", "Title", "Reference", "Duration"}
	for index, plannerOutput := range p.output {
		records[index+1] = plannerOutput.ToRecord()
	}

	return records
}
