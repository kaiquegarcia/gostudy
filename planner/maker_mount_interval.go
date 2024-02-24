package planner

import (
	"time"

	"github.com/kaiquegarcia/gostudy/v2/stream"
)

func (p *Maker) mountInterval(hgi *HourGradeInterval) error {
	initialStr := hgi.Start.Format(LayoutTimeOnly)
	endStr := hgi.End.Format(LayoutTimeOnly)
	p.logger.Debug("starting procedure for interval %s-%s", initialStr, endStr)
	var (
		previousDiscipline *Discipline
		previousSubject    string
		isFirst            = true
	)
	p.checkedDisciplinesCount = 0

	loopCounter := 0
	for {
		loopCounter++
		startStr := hgi.Start.Format(LayoutTimeOnly)
		p.logger.Debug("inner loop %d, checking if interval is still able to proceed", loopCounter)
		if !hgi.Start.Before(hgi.End) {
			p.logger.Debug("already reached the end of the interval, breaking at inner loop %d", loopCounter)
			p.checkedDisciplinesCount = 0
			break
		}

		p.logger.Debug(
			"inner loop %d still have time left (%s - %s), checking if already explored all possibilities",
			loopCounter, startStr, endStr,
		)
		if p.hasExploredAllDisciplines() {
			p.logger.Debug("already explored all possibilities, breaking at inner loop %d", loopCounter)
			p.checkedDisciplinesCount = 0
			break
		}

		p.logger.Debug("still have disciplines to explore, checking current discipline index %d", p.currentDisciplineIndex)
		discipline := p.disciplines[p.currentDisciplineIndex]
		p.logger.Debug("which means '%s' ~ checking if it's already finished", discipline.Name)
		if p.isDisciplineFinished(p.currentDisciplineIndex) {
			p.logger.Debug("discipline '%s' is already finished, adding gap only if current duration is higher than zero", discipline.Name)
			p.nextDiscipline(hgi)
			previousDiscipline = discipline
			continue
		}

		p.logger.Debug("discipline '%s' it's not finished, checking if it already exhausted daily limit (%s)", discipline.Name, discipline.DailyLimit)
		if p.currentDayDisciplineDuration >= discipline.DailyLimit {
			p.logger.Debug("discipline '%s' already exhausted daily limit, adding gap (if duration is higher than zero)", discipline.Name)
			p.nextDiscipline(hgi)
			previousDiscipline = discipline
			continue
		}

		p.logger.Debug("discipline '%s' did not exhaust daily limit, checking next content", discipline.Name)
		content, err := discipline.Next()
		if err == stream.ErrEOF {
			p.logger.Debug("discipline '%s' reached the end of content list, checking if all disciplines finished", discipline.Name)
			p.finishedDisciplinesIndexes = append(p.finishedDisciplinesIndexes, p.currentDisciplineIndex)

			if p.isAllDisciplinesFinished() {
				p.logger.Debug("all disciplines finished, ending planner mount! last discipline = %s", discipline.Name)
				p.outputWriter.Flush()
				return stream.ErrEOF
			}

			p.logger.Debug("still have disciplines to work on. getting next discipline, adding gap only if current duration is higher than zero")
			p.nextDiscipline(hgi)
			continue
		}

		if err != nil {
			p.logger.Error(err, "could not get next discipline content")
			return err
		}

		p.logger.Debug("discipline's content retrieved. checking if we should include a gap before the content")
		totalDuration := content.Duration
		var preGap time.Duration = 0
		if isFirst {
			p.logger.Debug("it's the first content of this time interval, no gap is required")
			preGap = 0
		} else if content.Subject != previousSubject && previousSubject != "" && previousDiscipline == discipline {
			totalDuration += discipline.SubjectGap
			preGap = discipline.SubjectGap
			p.logger.Debug("it's a new subject of the same discipline, adding gap of %s", preGap)
		} else if content.Subject == previousSubject {
			totalDuration += discipline.ContentGap
			preGap = discipline.ContentGap
			p.logger.Debug("it's a new content of the same subject, adding gap of %s", preGap)
		} else {
			p.logger.Debug("it's a new content from other disciplines, no gap is required")
		}

		p.logger.Debug("checking if discipline current duration (%s) + totalDuration (%s) exhaust discipline daily limit (%s)", p.currentDayDisciplineDuration, totalDuration, discipline.DailyLimit)
		if (p.currentDayDisciplineDuration + totalDuration) > discipline.DailyLimit {
			p.logger.Debug("discipline '%s's gap exhausts daily limit, getting next discipline, adding gap (if duration is higher than zero)", discipline.Name)
			p.nextDiscipline(hgi)
			previousDiscipline = discipline
			continue
		}

		p.logger.Debug("discipline '%s's gap doesn't exhaust daily limit", discipline.Name)

		p.logger.Debug("checking if the content + gap (%s) can be added to time interval (%s-%s)", totalDuration, startStr, endStr)
		if hgi.End.Sub(hgi.Start) < totalDuration {
			content.Attempts++
			p.logger.Debug(
				"total duration is higher than the time left, checking if this content attempts is higher than %d attempts",
				MaxContentAttemptsAllowed,
			)
			if content.Attempts > MaxContentAttemptsAllowed {
				return ErrContentDurationUnplayable
			}

			p.logger.Debug("content attempts is only %d, so we can attempt again next time", content.Attempts)
			err = discipline.Back()
			if err != nil {
				p.logger.Error(err, "could not step back on the discipline's content")
				return err
			}

			p.logger.Debug("as discipline '%s' can't fill with the current content, we'll call the next discipline", discipline.Name)
			p.nextDiscipline(hgi)
			continue
		}

		p.logger.Debug("content can be added to time interval, adding to PlannerOutput list")

		output := Output{
			Time:       hgi.Start.Add(preGap),
			Discipline: discipline,
			Content:    content,
		}

		err = p.outputWriter.Write(output.ToRecord())
		if err != nil {
			return err
		}

		previousSubject = content.Subject
		previousDiscipline = discipline
		p.currentDayDisciplineDuration += totalDuration
		hgi.Start = hgi.Start.Add(totalDuration)
		isFirst = false
		p.logger.Debug("inner loop %d finished, starting next", loopCounter)
	}

	p.outputWriter.Flush()

	p.logger.Debug("procedure for interval  %s-%s finished, calling next interval", initialStr, endStr)
	return nil
}
