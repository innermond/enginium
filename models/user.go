package models

type User struct {
	Id       int
	Username string
	password string
	salt     string
}

func (u *User) Get(id int) error {
	sql := "select id, username, password, salt from users where id = ?"
	err := DB.QueryRow(sql, id).Scan(&u.Id, &u.Username, &u.password, &u.salt)

	switch {
	case err != nil:
		return err
	default:
		return nil
	}
}

func (u *User) Authenticate(uname string, pwd string) int {
	var uid int = 0
	sql := "select id from users where username = ? and password = sha2(concat(?, salt), 256) limit 1"
	err := DB.QueryRow(sql, uname, pwd).Scan(&uid)
	if err != nil {
		return 0
	}
	return uid
}

func (u *User) String() string {
	return u.Username
}

func CreateUser(uname string, pwd string) int {
	// TODO: wrap into a transaction
	var uid int = 0
	sql := "insert into users (username, password) values(?, ?)"
	_, err := DB.Exec(sql, uname, pwd)
	if err != nil {
		return 0
	}
	sql = "select last_insert_id()"
	err = DB.QueryRow(sql).Scan(&uid)
	if err != nil {
		return 0
	}
	return uid
}
