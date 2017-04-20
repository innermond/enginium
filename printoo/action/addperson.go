package action

import (
	"github.com/innermond/printoo/printoo"
)

type PersonManager interface {
	AddPerson(printoo.Person) (printoo.Person, error)
	EditPerson(printoo.Person) error
	DeletePerson(printoo.Person) error
	GetPerson(int) (printoo.Person, error)
	AddExtraPhone(int, string) error
	DeleteExtraPhone(int, string) error
}

func NewHave(db *printoo.DB) PersonManager {
	return &personManager{db}
}

type personManager struct {
	*printoo.DB
}

func (my *personManager) AddPerson(p printoo.Person) (printoo.Person, error) {
	sql := "insert into persons (longname, phone, email, is_male, address, is_client, is_contractor) values(?, ?, ?, ?, ?, ? ,?)"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return p, err
	}
	defer stm.Close()

	res, err := stm.Exec(p.Longname, p.Phone, p.Email, p.IsMale, p.Address, p.IsClient, p.IsContractor)
	if err != nil {
		return p, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return p, err
	}
	p.Id = int(lid)
	return p, nil
}

func (my *personManager) EditPerson(p printoo.Person) error {
	sql := "update persons set longname=?, phone=?, email=?, is_male=?, address=?, is_client=?, is_contractor=? where id=?"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(p.Longname, p.Phone, p.Email, p.IsMale, p.Address, p.IsClient, p.IsContractor, p.Id)
	if err != nil {
		return err
	}

	return nil
}

func (my *personManager) DeletePerson(p printoo.Person) error {
	sql := "delete from persons where id=?"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(p.Id)
	if err != nil {
		return err
	}

	return nil
}

func (my *personManager) GetPerson(pid int) (printoo.Person, error) {
	var p printoo.Person
	sql := `select 
	id, longname, phone, email, 
	(is_male=true), -- to fit pointer boolean printoo.Person.IsMale
	address, is_client, is_contractor 
	from persons
	where id=?`
	err := my.DB.QueryRow(sql, pid).Scan(
		&p.Id,
		&p.Longname,
		&p.Phone,
		&p.Email,
		&p.IsMale,
		&p.Address,
		&p.IsClient,
		&p.IsContractor)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (my *personManager) AddExtraPhone(pid int, pp string) error {
	sql := "insert into person_phones (person_id, phone) values(?, ?)"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(pid, pp)
	if err != nil {
		return err
	}
	return nil
}

func (my *personManager) DeleteExtraPhone(pid int, pp string) error {
	sql := "delete from person_phones where person_id=? and phone=? limit 1"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(pid, pp)
	if err != nil {
		return err
	}

	return nil
}
