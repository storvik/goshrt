package http

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		s.InfoLog.Printf("%s %s used %s\n", r.Method, r.URL, time.Since(t))
	})
}

func (s *Server) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valid, err := s.Auth.Validate(r.Header.Get("Authorization"))
		if err != nil || !valid {
			if err != nil {
				s.ErrorLog.Printf("Could not validate token, %s\n", err.Error())
			}
			response, _ := json.Marshal(map[string]string{"response": "forbidden, could not authenticate"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write(response)
			return
		}
		next.ServeHTTP(w, r)
	})
}
