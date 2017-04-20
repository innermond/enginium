package services

import "github.com/innermond/printoo/models"

type User interface {
	Read(int) (*models.User, error)
	Authenticate(string, string) int
}

type CreateUserData struct {
	Username string
	Password string
}

type user struct{}

func NewUser() User {
	return &user{}
}

func (s *user) Read(id int) (*models.User, error) {
	User := &models.User{}
	err := User.Get(id)
	if err != nil {
		return nil, err
	}
	return User, nil
}

func (s *user) Authenticate(u string, p string) int {
	User := &models.User{}
	return User.Authenticate(u, p)
}

func CreateUser(ud CreateUserData) int {
	return models.CreateUser(ud.Username, ud.Password)
}
