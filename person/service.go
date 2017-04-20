package person

import (
	"errors"
	"log"

	"github.com/innermond/printoo/models"
	"github.com/innermond/printoo/services"
)

type Service interface {
	Read(int) (models.PersonModel, error)
	CreatePerson(Data) (Data, error)
}

type ErrDataValid error
type ErrCreate error

type service struct {
	model models.PersonModel
}

func NewService(m models.PersonModel) Service {
	return &service{m}
}

func (s *service) Read(id int) (models.PersonModel, error) {
	Person := s.model
	err := Person.Get(id)
	if err != nil {
		return nil, err
	}
	return Person, nil
}

func (s *service) CreatePerson(pd Data) (Data, error) {
	if ok := pd.Validate(); ok == false {
		log.Println(pd.Errors)
		return pd, ErrDataValid(errors.New("invalid data"))
	}

	var (
		err error
		pid int
	)
	s.model = pd.toPerson()
	p := s.model.Me()
	pid, err = s.model.CreatePerson()
	if err != nil {
		log.Println(err)
		return pd, ErrCreate(errors.New("cannot create"))
	}

	out := Data{*p, services.Errors(nil)}
	out.Id = pid

	return out, nil
}
