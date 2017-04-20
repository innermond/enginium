package handlers

import (
	"github.com/innermond/printoo/services"
	"log"
	"net/http"
)

type Token struct {
	services.Token
}

func (s *Token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		u := r.FormValue("username")
		p := r.FormValue("password")
		log.Println("auth credentials:", u, p)
		jwtString, err := s.Make(u, p)
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		w.Write([]byte(jwtString))
	default:
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}
}

func NewToken(s services.Token) *Token {
	return &Token{s}
}
