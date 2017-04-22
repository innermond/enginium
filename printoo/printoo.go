package printoo

import "fmt"

type Person struct {
	Id           int     `schema:"id" json:"id"`
	Longname     string  `schema:"longname" json:"longname"`
	Phone        *string `schema:"phone" json:"phone"`
	Email        *string `schema:"email" json:"email"`
	IsMale       *bool   `schema:"is_male" json:"is_male"`
	Address      string  `schema:"address" json:"address"`
	IsClient     bool    `schema:"is_client" json:"is_client"`
	IsContractor bool    `schema:"is_contractor" json:"is_contractor"`
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
