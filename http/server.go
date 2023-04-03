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

	// Index and shortener routes
	r.Get("/", indexHandler())
	r.Get("/{slug}", shrtHandler())

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Should authenticate here r.Use()
		r.Route("/shrt", func(r chi.Router) {
			r.Post("/", shrtCreateHandler()) // POST          /shrt         - Create new shrt
			r.Route("/{slug}", func(r chi.Router) {
				r.Get("/", shrtGetHandler())       // GET     /shrt/{slug}  - Get shrt details
				r.Delete("/", shrtDeleteHandler()) // DELETE  /shrt/{slug}  - Delete shrt
			})
		})
	})

	return &Server{
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

func shrtHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler should lookup slug in db and redirect
	})
}

func shrtCreateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func shrtGetHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func shrtDeleteHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func shrtListHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
