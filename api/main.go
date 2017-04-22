package main

import (
	"log"
	"net/http"
)

func main() {

	api := NewApi()
	defer api.Close()

	http.Handle("/person", api.Person)

	certPath := "server.pem"
	keyPath := "server.key"
	addr := "localhost:3000"
	log.Println("Start server " + addr)
	log.Fatal(http.ListenAndServeTLS(addr, certPath, keyPath, nil))

}
