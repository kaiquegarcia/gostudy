package gostudy

import "time"

type PlannerOutput struct {
	Time       time.Time
	Discipline *Discipline
	Subject    *Subject
	Content    *Content
}

func (po PlannerOutput) ToRecord() []string {
	return []string{
		po.Time.Format(time.RFC3339),
		po.Discipline.Name,
		po.Subject.Name,
		po.Content.Title,
		po.Content.Reference,
		po.Content.Duration.String(),
	}
}
