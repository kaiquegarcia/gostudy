package gostudy

import "time"

type Subject struct {
	Name     string
	contents []*Content
	pointer  int
}

func NewSubject(name string) *Subject {
	return &Subject{
		Name:     name,
		contents: make([]*Content, 0),
		pointer:  -1,
	}
}

func (s *Subject) Add(c *Content) {
	s.contents = append(s.contents, c)
}

func (s *Subject) SetContents(cs []*Content) {
	s.contents = cs
	s.pointer = -1
}

func (s *Subject) Next() (*Content, error) {
	if s.pointer >= (len(s.contents) - 1) {
		return nil, ErrEndOfList
	}

	s.pointer++
	return s.contents[s.pointer], nil
}

func (s *Subject) Back() error {
	if s.pointer < 0 {
		return ErrBeginningOfList
	}

	s.pointer--
	return nil
}

func (s *Subject) TotalDuration() time.Duration {
	t := time.Duration(0)
	for _, c := range s.contents {
		t += c.Duration
	}

	return t
}
