package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/innermond/printoo/domino"
	"github.com/innermond/printoo/handlers"
	"github.com/innermond/printoo/person"
	"github.com/innermond/printoo/services"
	"github.com/joho/godotenv"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func init() {
	// load config values
	err := godotenv.Load()
	fatal(err)
	// build DNS string for database
	dbname := os.Getenv("DB_NAME")
	dbpwd := os.Getenv("DB_PWD")
	dns := fmt.Sprintf("root:%s@tcp(:3306)/%s", dbpwd, dbname)
	// connect to database
	err = services.Storage.Init(dns)
	fatal(err)
}

type api struct {
	Hello,
	User,
	Token,
	Person http.Handler
}

func (a *api) Clean() {
	services.Storage.End()
	log.Println("api cleaned")
}

func NewApi() *api {
	user := services.NewUser()
	userHandler := handlers.NewUser(user)

	token := services.NewToken(user)
	tokenHandler := handlers.NewToken(token)

	personService := person.NewService()
	personHandler := domino.Pieces(
		handlers.Note,
		/*handlers.Recover,*/
		//handlers.CheckToken,
		person.Handle(personService),
	).Roll(handlers.ConvertJson())
	log.Println("New api delivered")
	return &api{
		User:   userHandler,
		Token:  tokenHandler,
		Person: personHandler,
	}
}
