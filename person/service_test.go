package person_test

import (
	"testing"

	"github.com/innermond/printoo/models"
	"github.com/innermond/printoo/person"
)

func TestTest(t *testing.T) {
	t.Log("passing")
}

type fakePersonModel struct{ models.Person }

func (k *fakePersonModel) Get(id int) error {
	return nil
}

func (k *fakePersonModel) CreatePerson() (int, error) {
	return k.Id, nil
}

func TestCreatePerson(t *testing.T) {
	m := &fakePersonModel{}
	s := person.NewService(m)
	pd := person.NewData(map[string][]string{
		"id": {"0"}, "longname": {"aaa"},
	})
	dt, err := s.CreatePerson(pd)
	t.Log(dt, err)
}
