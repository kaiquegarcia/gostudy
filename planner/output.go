package planner

import "time"

type Output struct {
	Time       time.Time
	Discipline *Discipline
	Content    *Content
}

func (po Output) ToRecord() []string {
	return []string{
		po.Time.Format(time.RFC3339),
		po.Discipline.Name,
		po.Content.Subject,
		po.Content.Title,
		po.Content.Reference,
		po.Content.Duration.String(),
	}
}
