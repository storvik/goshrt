package http

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/storvik/goshrt"
)

// Server type, must be global in order to addach interfaces used
// in http routes.
type Server struct {
	ln       net.Listener
	server   *http.Server
	router   *chi.Mux
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	// Interfaces required in various http routes etc
	Auth      goshrt.Authorizer
	ShrtStore goshrt.ShrtStorer
}

func NewServer(l *log.Logger, p string) *Server {
	// Create router
	r := chi.NewRouter()

	// Create server instance and attach router
	s := &Server{
		server: &http.Server{
			Addr:         p,
			Handler:      r,
			ErrorLog:     l,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 10,
			IdleTimeout:  time.Second * 60,
		},
		router:   r,
		ErrorLog: l,
	}

	// Index and shortener routes
	r.Get("/", indexHandler())
	r.Get("/{slug}", s.shrtHandler())

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Should authenticate here r.Use()
		r.Route("/shrt", func(r chi.Router) {
			r.Post("/", s.shrtCreateHandler()) // POST          /shrt         - Create new shrt
			r.Route("/{slug}", func(r chi.Router) {
				r.Get("/", s.shrtGetHandler())       // GET     /shrt/{slug}  - Get shrt details
				r.Delete("/", s.shrtDeleteHandler()) // DELETE  /shrt/{slug}  - Delete shrt
			})
		})
	})

	return s

}

func (s *Server) ListenAndServe() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.server.Shutdown(ctx)
	return nil
}

func indexHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(map[string]string{"response": "Index endpoint. Probably should do something clever here."})
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

func (s *Server) shrtHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler should lookup slug in db and redirect
	})
}

func (s *Server) shrtCreateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shrt := new(goshrt.Shrt)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&shrt)
		if err != nil {
			s.ErrorLog.Printf("Could not decode json: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If slug is empty, generate random slug
		if shrt.Slug == "" {
			// TODO: Make slug length configurable
			shrt.Slug = goshrt.GenerateSlug(7)
		}

		err = s.ShrtStore.CreateShrt(shrt)
		if err != nil {
			response, _ := json.Marshal(map[string]string{"response": "error storing shrt"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			s.ErrorLog.Printf("Could not store shrt to database: %s\n", err)
			return
		}

		response, _ := json.Marshal(shrt)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	})
}

func (s *Server) shrtGetHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shrt := new(goshrt.Shrt)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&shrt)
		if err != nil {
			s.ErrorLog.Printf("Could not decode json: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if shrt.ID > 0 {
			// Use ID if ID is set
			shrt, err = s.ShrtStore.ShrtByID(shrt.ID)
		} else {
			shrt, err = s.ShrtStore.Shrt(shrt.Domain, shrt.Slug)
		}
		if err != nil {
			response, _ := json.Marshal(map[string]string{"response": "error retrieving"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			s.ErrorLog.Printf("Could not get shrt from database: %s\n", err)
			return
		}

		response, _ := json.Marshal(shrt)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}

func (s *Server) shrtDeleteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (s *Server) shrtListHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
