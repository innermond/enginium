package person

import (
	"encoding/json"
	"log"
	"net/http"
)

/*func Handle(s Service) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		f := func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				/*r.ParseForm()
				d := r.PostForm
				p, ok, verr := printoo.NewPerson(d)
				if !ok {
					//t.Fatalf("%v", verr)
				}
				p, err := do.AddPerson(p)
				if err != nil {
					//t.Fatalf("insert: %v", err)
				}
				if err != nil || pd.HasErrors() {
					status := 500
					if _, ok := err.(ErrDataValid); ok {
						status = 412 // precondition failed ~ validation failed
					}
					w.WriteHeader(status)
				}
				ctx := services.ContextedMistaker(r.Context(), pd)
				r = r.WithContext(ctx)
				h.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(f)
	}

}*/

func ServicedHandler(s *Service) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			r.ParseForm()
			d := r.PostForm
			p, err := s.AddPerson(d)
			if err != nil {
				log.Printf("%v \n", err)
				w.WriteHeader(412)
				json.NewEncoder(w).Encode(err)
				return
			}
			jsn, _ := json.Marshal(p)
			w.Write(jsn)
		}
	}
	return http.HandlerFunc(f)

}
