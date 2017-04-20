package person

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"github.com/gorilla/schema"
	"github.com/innermond/printoo/models"
	"github.com/innermond/printoo/services"
)

func isPhone(p string, l int) bool {
	re := regexp.MustCompile("[0-9]+")
	nb := re.FindAllString(p, -1)
	s := []string{}
	for _, n := range nb {
		s = append(s, n)
	}
	phone := strings.Join(s, "")
	out := len(phone) >= l
	return out
}

func minmax(min, max, val int, eq ...bool) bool {
	is := val < min || val > max
	if len(eq) > 0 && eq[0] {
		is = val <= min || val >= max

	}
	return is
}

type Data struct {
	models.Person
	services.Errors
}

func (pd *Data) Ok(fn string) {
	switch fn {
	case "longname":
		longname := strings.TrimSpace(pd.Longname)
		printable := true
		for _, ch := range longname {
			printable = unicode.IsPrint(ch)
			if !printable {
				break
			}
		}
		minLongname := 4
		maxLongname := 50
		l := len(longname)
		if longname == "" ||
			minmax(minLongname, maxLongname, l) ||
			!printable {
			pd.ErrMessage("longname", longname)
		}
	case "phone":
		minPhone := 10
		maxPhone := 15
		phone := strings.TrimSpace(pd.Phone)
		l := len(phone)
		if minmax(minPhone, maxPhone, l) ||
			!isPhone(phone, 10) {
			pd.ErrMessage("phone", phone)
		}
	}
}

func (pd *Data) Validate() bool {
	for _, fn := range []string{"longname", "phone"} {
		pd.Ok(fn)
	}
	return !pd.HasErrors()
}

func (pd Data) toPerson() *models.Person {
	return &models.Person{
		pd.Id, pd.Longname, pd.Phone, pd.Email, pd.IsMale, pd.Address, pd.IsClient, pd.IsContractor,
	}
}

func (pd Data) Errorless() interface{} {
	return interface{}(pd.toPerson())
}

func (pd Data) HasErrors() bool {
	return len(pd.Errors) > 0
}

func (pd Data) GetErrors() services.Errors {
	return pd.Errors
}

func NewData(d map[string][]string) Data {

	var pd Data
	pd.Errors = services.Errors{}

	decoder := schema.NewDecoder()
	err := decoder.Decode(&pd, d)
	if err != nil {
		log.Println(err)
		pd.ErrMessage("decode", err)
		return pd
	}

	return pd
}
