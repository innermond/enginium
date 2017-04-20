package models

import (
	"errors"
	"fmt"
	"log"
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

type PersonModel interface {
	Get(int) error
	CreatePerson() (int, error)
	Me() *Person
}

func (u *Person) Get(id int) error {
	sql := "select id, personname, password, salt from persons where id = ?"
	err := DB.QueryRow(sql, id).Scan(&u.Id, &u.Longname)

	switch {
	case err != nil:
		return err
	default:
		return nil
	}
}

func (p *Person) String() string {
	return p.Longname
}

func (p *Person) Me() *Person {
	return p
}

func (p *Person) CreatePerson() (int, error) {

	if p.Id != 0 {
		return 0, errors.New(fmt.Sprintf("person id %q should be zero", p.Id))
	}
	tx, err := DB.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sql := "insert into persons (longname, phone, email, is_male, address, is_client, is_contractor) values(?, nullif(?,''), nullif(?, ''), ?, nullif(?, ''), ? ,?)"
	stm, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stm.Close()

	_, err = stm.Exec(p.Longname, p.Phone, p.Email, p.IsMale, p.Address, p.IsClient, p.IsContractor)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	sql = "select last_insert_id()"
	stm, _ = tx.Prepare(sql)
	err = stm.QueryRow().Scan(&p.Id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	log.Println("models p.Id: ", p.Id)
	return p.Id, nil
}
