package gostudy

import "time"

func (p *Planner) mountDate(date time.Time) error {
	dateStr := date.Format(LayoutDateOnly)
	p.logger.Debug("retrieving time intervals for date %s", dateStr)
	intervals, err := p.hg.IntervalsFor(date)
	if err != nil {
		p.logger.Error(err, "could not retrieve time intervals for date %s", dateStr)
		return err
	}

	p.currentDayDisciplineDuration = 0
	p.logger.Debug("%d intervals found, start loop", len(intervals))
	for _, hgi := range intervals {
		err = p.mountInterval(hgi)
		if err != nil {
			return err
		}
	}

	return nil
}
