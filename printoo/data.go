package printoo

import (
	"strings"
	"unicode"

	"github.com/gorilla/schema"
	"github.com/innermond/printoo/printoo/validation"
)

func NewPerson(d map[string][]string) (Person, bool, validation.Errors) {
	var pd Person
	ok := true
	ve := make(validation.Errors)

	decoder := schema.NewDecoder()
	err := decoder.Decode(&pd, d)
	if err != nil {
		ve.Err("decode", err)
		return pd, false, ve
	}
	ok, ve = ValidateNewPerson(pd)
	return pd, ok, ve
}

func NewPersonExtras(d map[string][]string) (Person, []ExtraPhone, []ExtraEmail, bool, validation.Errors) {
	var (
		xp []ExtraPhone
		xm []ExtraEmail
	)
	p, ok, verr := NewPerson(d)
	if !ok {
		return p, xp, xm, false, verr
	}
	if pp, ok := d["phone"]; ok && len(pp) > 1 {
		pp = pp[:len(pp)-1]
		for _, ph := range pp {
			if validation.IsPhone(ph, 10) {
				xp = append(xp, ExtraPhone(ph))
			} else {
				verr.Err("extraphone", ph)
			}
		}
	}
	// extra emails
	if mm, ok := d["email"]; ok && len(mm) > 1 {
		mm = mm[:len(mm)-1]
		for _, em := range mm {
			if validation.IsEmail(em) {
				xm = append(xm, ExtraEmail(em))
			} else {
				verr.Err("extraemail", em)
			}
		}
	}

	if len(verr) > 0 {
		return p, xp, xm, false, verr
	}

	return p, xp, xm, true, verr
}

func validatePerson(p Person) (bool, validation.Errors) {
	ok := true
	ve := make(validation.Errors)

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

	if p.Phone != nil {
		minPhone := 10
		maxPhone := 15
		phone := strings.TrimSpace(*p.Phone)
		if phone != "" {
			l = len(phone)
			if validation.MinMax(minPhone, maxPhone, l) ||
				!validation.IsPhone(phone, 10) {
				ve.Err("phone", phone)
				ok = false
			}
		}
	}

	return ok, ve
}

func ValidateNewPerson(p Person) (bool, validation.Errors) {
	ok := true
	ve := make(validation.Errors)

	if p.Id != 0 {
		ve.Err("id", "not zero")
		return false, ve
	}
	ok, ve = validatePerson(p)
	return ok, ve
}

func ValidateEditPerson(p Person) (bool, validation.Errors) {
	ok := true
	ve := make(validation.Errors)

	if p.Id == 0 {
		ve.Err("id", "is zero")
		return false, ve
	}
	ok, ve = validatePerson(p)
	return ok, ve
}
