package handlers

import (
	"fmt"
	"github.com/innermond/printoo/services"
	"net/http"
	"strings"
)

type User struct {
	services.User
}

func (s *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		segments := strings.Split(r.URL.Path, "/")
		fmt.Println(segments)
		u, err := s.Read(1)
		if err != nil {
			InternalServerError(w)
			return
		}
		w.Write([]byte(u.String()))
	case "POST":
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w)
			return
		}
		ud := services.CreateUserData{
			r.FormValue("username"),
			r.FormValue("password"),
		}
		uid := services.CreateUser(ud)
		if uid == 0 {
			InternalServerError(w)
			return
		}
		w.Write([]byte(string(uid)))
	default:
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
	}
}

func NewUser(s services.User) *User {
	return &User{s}
}

func InternalServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}
