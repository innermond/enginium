package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/innermond/printoo/person"
	"github.com/innermond/printoo/printoo"
	"github.com/innermond/printoo/printoo/action"
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
}

type api struct {
	Person http.Handler
	Close  func()
}

func NewApi() *api {
	// build DNS string for database
	dbname := os.Getenv("DB_NAME")
	dbpwd := os.Getenv("DB_PWD")
	dns := fmt.Sprintf("root:%s@tcp(:3306)/%s", dbpwd, dbname)
	// init api
	db, err := printoo.Open(dns)
	fatal(err)

	do := action.NewHave(db)
	if do == nil {
		panic("no do")
	}
	personService := person.NewService(do)
	personHandler := person.ServicedHandler(personService)

	log.Println("New api delivered")

	return &api{
		Person: personHandler,
		Close: func() {
			log.Println("api closing...")
			db.Close()
		},
	}
}
