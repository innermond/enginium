package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/innermond/printoo/api"
	"github.com/innermond/printoo/printoo"
	"github.com/joho/godotenv"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	confPath := flag.String("conf", "./conf.json", "conf file path relative to main")
	conf, err := printoo.ConfigFrom(*confPath)
	fatal(err)

	err = godotenv.Load(conf.Env)
	fatal(err)

	a := api.NewApi()
	defer a.Close()

	http.Handle("/person", a.Person)

	certPath := conf.ServerPem
	keyPath := conf.ServerKey
	addr := conf.ServerAddr

	log.Println("Start server " + addr)
	log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, nil))
}
