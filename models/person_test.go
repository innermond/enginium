package models

import "testing"

func TestPersonString(t *testing.T) {

	t.Log("person model String method returns Longname field value")
	uses := []struct {
		p    Person
		want string
	}{
		{Person{}, ""},
		{Person{Id: 1, Longname: "gb"}, "gb"},
		{Person{Id: 2, Longname: "Gabriel Braila"}, "Gabriel Braila"},
	}
	for _, use := range uses {
		if use.p.String() != use.want {
			t.Errorf("Expected %q got %q", use.want, use.p.String())
		}
	}
}

func TestPersonCreate(t *testing.T) {
	InitDatabase("root:M0b1d1c3@tcp(:3306)/printoo")
	t.Log("create from an empty person struct")
	p := &Person{}
	pid, err := p.CreatePerson()
	if err != nil || pid == 0 || p.Id != pid {
		t.Errorf("Expected non error %s got id %d", err.Error(), pid)
	}
}
