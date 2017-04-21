package printoo

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gorilla/schema"
	"github.com/innermond/printoo/printoo/validation"
)

type Person struct {
	Id           int    `schema:"id" json:"id"`
	Longname     string `schema:"longname" json:"longname"`
	Phone        string `schema:"phone" json:"phone"`
	Email        string `schema:"email" json:"email"`
	IsMale       *bool  `schema:"is_male" json:"is_male"`
	Address      string `schema:"address" json:"address"`
	IsClient     bool   `schema:"is_client" json:"is_client"`
	IsContractor bool   `schema:"is_contractor" json:"is_contractor"`
}

func NewPerson(d map[string][]string) (Person, bool, ValidationErrors) {
	var pd Person
	ok := true
	ve := make(ValidationErrors)

	decoder := schema.NewDecoder()
	err := decoder.Decode(&pd, d)
	if err != nil {
		ve.Err("decode", err)
		return pd, false, ve
	}

	ok, ve = ValidateNewPerson(pd)
	return pd, ok, ve
}

type ValidationErrors map[string][]string

func (e ValidationErrors) Err(fld string, val interface{}) {
	qv := "%q"
	if err, ok := val.(error); ok == true {
		val = err.Error()
		qv = "%v"
	}
	e[fld] = append(e[fld], fmt.Sprintf(qv, val))
}

func validatePerson(p Person) (bool, ValidationErrors) {
	ok := true
	ve := make(ValidationErrors)

	if p.Longname == "" {
		ve.Err("longname", "required")
		return false, ve
	}

	longname := strings.TrimSpace(p.Longname)
	printable := true
	for _, ch := range longname {
		printable = unicode.IsPrint(ch)
		if !printable {
			ok = false
			ve.Err("longname", "forbidden character")
			break
		}
	}

	minLongname := 4
	maxLongname := 50
	l := len(longname)
	if longname == "" ||
		validation.MinMax(minLongname, maxLongname, l) ||
		!printable {
		ve.Err("longname", longname)
		ok = false
	}

	minPhone := 10
	maxPhone := 15
	phone := strings.TrimSpace(p.Phone)
	if phone != "" {
		l = len(phone)
		if validation.MinMax(minPhone, maxPhone, l) ||
			!validation.IsPhone(phone, 10) {
			ve.Err("phone", phone)
			ok = false
		}
	}

	return ok, ve

}

func ValidateNewPerson(p Person) (bool, ValidationErrors) {
	ok := true
	ve := make(ValidationErrors)

	if p.Id != 0 {
		ve.Err("id", "not zero")
		return false, ve
	}
	ok, ve = validatePerson(p)
	return ok, ve
}

func ValidateEditPerson(p Person) (bool, ValidationErrors) {
	ok := true
	ve := make(ValidationErrors)

	if p.Id == 0 {
		ve.Err("id", "is zero")
		return false, ve
	}
	ok, ve = validatePerson(p)
	return ok, ve
}

type ExtraPhone string

func (ep *ExtraPhone) Scan(src interface{}) error {
	if extra, ok := src.([]byte); ok {
		*ep = ExtraPhone(extra)
		return nil
	}
	return fmt.Errorf("cannot %T %v", ep, src)
}

type ExtraEmail string

func (em *ExtraEmail) Scan(src interface{}) error {
	if extra, ok := src.([]byte); ok {
		*em = ExtraEmail(extra)
		return nil
	}
	return fmt.Errorf("cannot %T %v", em, src)
}
