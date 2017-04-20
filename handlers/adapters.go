package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/innermond/printoo/services"
)

func ConvertJson() http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {
		pd := r.Context().Value(services.ConvertJsonKey)
		jsn, _ := json.Marshal(pd)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsn)
	}

	return http.HandlerFunc(f)
}

func CheckToken(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		log.Println("checktoken")
		l := r.Header.Get("Unlock")
		log.Println(l)
		if l != "" {
			ctx := context.WithValue(r.Context(), "token", l)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(500)
			w.Write([]byte("no unlock header"))
		}
	}

	return http.HandlerFunc(f)
}

func Recover(h http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func Note(h http.Handler) http.Handler {

	f := func(w http.ResponseWriter, r *http.Request) {

		t0 := time.Now()
		h.ServeHTTP(w, r)
		t8 := time.Now()
		log.Printf("%s %q %v\n", r.Method, r.URL.String(), t8.Sub(t0))
	}

	return http.HandlerFunc(f)
}
