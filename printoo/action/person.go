package action

import "github.com/innermond/printoo/printoo"

type PersonManager interface {
	AddPerson(printoo.Person) (printoo.Person, error)
	AddPersonWithExtras(printoo.Person, []printoo.ExtraPhone, []printoo.ExtraEmail) (printoo.Person, error)

	EditPerson(printoo.Person) error
	DeletePerson(printoo.Person) error
	GetPerson(int) (printoo.Person, error)

	ExtraPhoneManager
	ExtraEmailManager
}

type ExtraPhoneManager interface {
	AddExtraPhone(int, string) error
	DeleteExtraPhone(int, string) error
	GetExtraPhones(int) ([]printoo.ExtraPhone, error)
}

type ExtraEmailManager interface {
	AddExtraEmail(int, string) error
	DeleteExtraEmail(int, string) error
	GetExtraEmails(int) ([]printoo.ExtraEmail, error)
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

func (my *personManager) AddPersonWithExtras(p printoo.Person, pp []printoo.ExtraPhone, mm []printoo.ExtraEmail) (printoo.Person, error) {
	tx, err := my.DB.Begin()
	if err != nil {
		return p, err
	}
	sql := "insert into persons (longname, phone, email, is_male, address, is_client, is_contractor) values(?, ?, ?, ?, ?, ? ,?)"
	stmp, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return p, err
	}
	defer stmp.Close()

	res, err := stmp.Exec(p.Longname, p.Phone, p.Email, p.IsMale, p.Address, p.IsClient, p.IsContractor)
	if err != nil {
		tx.Rollback()
		return p, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return p, err
	}
	p.Id = int(lid)

	sql = "insert into person_phones (person_id, phone) values(?, ?)"
	stmph, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return p, err
	}
	defer stmph.Close()
	for _, ph := range pp {
		_, err = stmph.Exec(p.Id, string(ph))
		if err != nil {
			tx.Rollback()
			return p, err
		}
	}

	sql = "insert into person_emails (person_id, email) values(?, ?)"
	stmpe, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return p, err
	}
	defer stmpe.Close()
	for _, em := range mm {
		_, err = stmpe.Exec(p.Id, string(em))
		if err != nil {
			tx.Rollback()
			return p, err
		}
	}
	tx.Commit()
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

func (my *personManager) GetExtraPhones(pid int) ([]printoo.ExtraPhone, error) {
	var phones []printoo.ExtraPhone
	sql := "select phone from person_phones where person_id = ?"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return phones, err
	}
	defer stm.Close()
	rows, err := stm.Query(pid)
	if err != nil {
		return phones, err
	}
	defer rows.Close()
	for rows.Next() {
		var phone printoo.ExtraPhone
		err = rows.Scan(&phone)
		if err != nil {
			return phones, err
		}
		phones = append(phones, phone)
	}
	err = rows.Err()
	if err != nil {
		return phones, err
	}
	return phones, nil
}

func (my *personManager) AddExtraEmail(pid int, em string) error {
	sql := "insert into person_emails (person_id, email) values(?, ?)"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(pid, em)
	if err != nil {
		return err
	}
	return nil
}

func (my *personManager) DeleteExtraEmail(pid int, em string) error {
	sql := "delete from person_emails where person_id=? and email=? limit 1"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(pid, em)
	if err != nil {
		return err
	}

	return nil
}

func (my *personManager) GetExtraEmails(pid int) ([]printoo.ExtraEmail, error) {
	var emails []printoo.ExtraEmail
	sql := "select email from person_emails where person_id = ?"
	stm, err := my.DB.Prepare(sql)
	if err != nil {
		return emails, err
	}
	defer stm.Close()
	rows, err := stm.Query(pid)
	if err != nil {
		return emails, err
	}
	defer rows.Close()
	for rows.Next() {
		var email printoo.ExtraEmail
		err = rows.Scan(&email)
		if err != nil {
			return emails, err
		}
		emails = append(emails, email)
	}
	err = rows.Err()
	if err != nil {
		return emails, err
	}
	return emails, nil
}
