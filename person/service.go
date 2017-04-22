package person

import (
	"github.com/innermond/printoo/printoo"
	"github.com/innermond/printoo/printoo/action"
)

type Service struct {
	m action.PersonManager
}

func NewService(m action.PersonManager) *Service {
	return &Service{m}
}

func (s *Service) AddPerson(d map[string][]string) (printoo.Person, error) {
	p, ok, verr := printoo.NewPerson(d)
	if !ok {
		return p, verr
	}
	p, err := s.m.AddPerson(p)
	if err != nil {
		return p, err
	}

	return p, nil
}
