package http

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/storvik/goshrt"
)

// TODO: Add tests using httptest

// Server type, must be global in order to addach interfaces used
// in http routes.
type Server struct {
	ln         net.Listener
	server     *http.Server
	router     *chi.Mux
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	SlugLength uint64

	// Interfaces required in various http routes etc
	Auth      goshrt.Authorizer
	ShrtStore goshrt.ShrtStorer
}

// NewServer creates new http server with given errorlogger and port.
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
	r.Get("/*", s.shrtHandler())

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(s.authorize)
		r.Use(s.requestLogger)
		r.Route("/shrt", func(r chi.Router) {
			r.Post("/", s.shrtCreateHandler()) // POST              /shrt
			r.Route("/{id_domain}", func(r chi.Router) {
				r.Get("/", s.shrtGetHandler())       // GET         /shrt/{id}
				r.Delete("/", s.shrtDeleteHandler()) // DELETE      /shrt/{id}
				r.Route("/{slug}", func(r chi.Router) {
					r.Get("/", s.shrtGetHandler()) // GET           /shrt/{domain}/{slug}
				})
			})
		})
		r.Route("/shrts", func(r chi.Router) {
			r.Get("/", s.shrtListHandler()) // GET                  /shrts
			r.Route("/{domain}", func(r chi.Router) {
				r.Get("/", s.shrtListHandler()) // GET              /shrts/{domain}
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

// TODO: Better looking error site
func (s *Server) shrtHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		slug := chi.URLParam(r, "*")
		if slug == "" {
			s.ErrorLog.Println("Could not get empty slug")
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Parse URL to get domain
		shrt, err := s.ShrtStore.Shrt(r.Host, slug)
		if err == goshrt.ErrNotFound {
			s.InfoLog.Println("Could not find, " + r.Host + "/" + slug)
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			s.ErrorLog.Println("Could not get shrt, " + err.Error())
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.InfoLog.Printf("%s %s used %s  --> %s\n", r.Method, r.URL, time.Since(t), shrt.Dest)
		http.Redirect(w, r, shrt.Dest, http.StatusMovedPermanently)
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
			shrt.Slug = goshrt.GenerateSlug(s.SlugLength)
		}
		if !shrt.ValidDest() || !goshrt.ValidateSlug(shrt.Slug) {
			response, _ := json.Marshal(map[string]string{"response": "error storing shrt"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			s.ErrorLog.Printf("Invalid request, destination address or slug is not valid\n")
			return
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
		var err error
		shrt := new(goshrt.Shrt)

		// If slug is set, domain shold be present
		// Else ID should be firsts
		if slug := chi.URLParam(r, "slug"); slug != "" {
			domain := chi.URLParam(r, "id_domain")
			shrt, err = s.ShrtStore.Shrt(domain, slug)
		} else {
			id := chi.URLParam(r, "id_domain")
			idInt, _ := strconv.Atoi(id)
			shrt, err = s.ShrtStore.ShrtByID(idInt)
		}

		if err == goshrt.ErrNotFound {
			response, _ := json.Marshal(map[string]string{"response": "error retrieving"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write(response)
			s.ErrorLog.Printf("Could not get shrt from database: %s\n", err)
			return
		} else if err != nil {
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
		id := chi.URLParam(r, "id_domain")
		idInt, _ := strconv.Atoi(id)
		shrt, err := s.ShrtStore.DeleteByID(idInt)
		if err == goshrt.ErrNotFound {
			response, _ := json.Marshal(map[string]string{"response": "could not find item with given id"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write(response)
			s.ErrorLog.Printf("Could not find and delete shrt with id %d, %s\n", idInt, err)
			return
		} else if err != nil {
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

func (s *Server) shrtListHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var shrts []*goshrt.Shrt

		domain := chi.URLParam(r, "domain")
		if domain == "" {
			shrts, err = s.ShrtStore.Shrts()
		} else {
			shrts, err = s.ShrtStore.ShrtsByDomain(domain)
		}

		if err == goshrt.ErrNotFound {
			response, _ := json.Marshal(map[string]string{"response": "error retrieving"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write(response)
			s.ErrorLog.Printf("Could not get shrts from database: %s\n", err)
			return
		} else if err != nil {
			response, _ := json.Marshal(map[string]string{"response": "error retrieving"})
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			s.ErrorLog.Printf("Could not get shrts from database: %s\n", err)
			return
		}

		response, _ := json.Marshal(shrts)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
}
