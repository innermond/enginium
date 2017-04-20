package person

import (
	"net/http"

	"github.com/innermond/printoo/services"
)

func Handle(s Service) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		f := func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				r.ParseForm()
				pd := NewData(r.PostForm)
				pd, err := s.CreatePerson(pd)
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

}
