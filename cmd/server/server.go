// fhdata - Far Horizons Data
//
// Copyright (c) 2022 Michael D Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package main

import (
	"bytes"
	"github.com/mdhender/fhdata"
	"github.com/mdhender/fhdata/internal/way"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func NewServer(host, port string, opts ...func(*Server) error) (*Server, error) {
	s := &Server{
		router: way.NewRouter(),
	}
	s.Addr = net.JoinHostPort(host, port)
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20 // 1mb?

	// apply the list of options to Store
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	s.router.HandleFunc("GET", "/manifest.json", s.manifestJsonV3)
	s.router.HandleFunc("GET", "/home", s.getHome())
	s.router.HandleFunc("GET", "/planets", s.getPlanets())
	s.router.HandleFunc("GET", "/planet/:id", s.getPlanet())
	s.router.HandleFunc("GET", "/species", s.getSpecies())
	s.router.HandleFunc("GET", "/specie/:id", s.getSpecie())
	s.router.HandleFunc("GET", "/systems", s.getSystems())
	s.router.HandleFunc("GET", "/system/:id", s.getSystem())
	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	return s, nil
}

type Server struct {
	http.Server
	router    *way.Router
	templates string // path to templates directory
	data      *fhdata.Cluster
}

type Option func(*Server) error

// Options turns a list of Option instances into an Option.
func Options(opts ...Option) Option {
	return func(s *Server) error {
		for _, opt := range opts {
			if err := opt(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithStore(store *fhdata.Cluster) Option {
	return func(s *Server) (err error) {
		s.data = store
		return nil
	}
}

func WithTemplates(root string) Option {
	return func(s *Server) (err error) {
		s.templates = filepath.Clean(root)
		return nil
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) getHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getHome: %s %s\n", r.Method, r.URL.Path)
		b, err := s.render("home", s.data)
		if err != nil {
			log.Printf("getHome: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getPlanet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getPlanet: %s %s\n", r.Method, r.URL.Path)
		id, err := strconv.Atoi(way.Param(r.Context(), "id"))
		if err != nil || !(0 < id && id <= len(s.data.Planets)) {
			//log.Printf("getPlanet: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		planet := s.data.Planets[id-1]
		b, err := s.render("planet", planet)
		if err != nil {
			log.Printf("getPlanet: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getPlanets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getPlanets: %s %s\n", r.Method, r.URL.Path)
		b, err := s.render("planets", s.data)
		if err != nil {
			log.Printf("getPlanets: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getSpecie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getSpecie: %s %s\n", r.Method, r.URL.Path)
		id, err := strconv.Atoi(way.Param(r.Context(), "id"))
		if err != nil || !(0 < id && id <= len(s.data.Species)) {
			log.Printf("getSpecie: %s %s: %+v\n", r.Method, r.URL.Path, err)
			log.Printf("getSpecie: %s %s: len(s.data.Species) %d id %q\n", r.Method, r.URL.Path, len(s.data.Species), way.Param(r.Context(), "id"))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		specie := s.data.Species[id-1]
		b, err := s.render("specie", specie)
		if err != nil {
			log.Printf("getSpecie: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getSpecies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getSpecies: %s %s\n", r.Method, r.URL.Path)
		b, err := s.render("species", s.data)
		if err != nil {
			log.Printf("getSpecies: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getSystem: %s %s\n", r.Method, r.URL.Path)
		id, err := strconv.Atoi(way.Param(r.Context(), "id"))
		if err != nil || !(0 < id && id <= len(s.data.Systems)) {
			//log.Printf("getSystem: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		system := s.data.Systems[id-1]
		b, err := s.render("system", system)
		if err != nil {
			log.Printf("getSystem: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) getSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("getSystems: %s %s\n", r.Method, r.URL.Path)
		b, err := s.render("systems", s.data)
		if err != nil {
			log.Printf("getSystems: %s %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

func (s *Server) manifestJsonV3(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"manifest_version":3,"name":"My Extension","version":"versionString"}`))
}

func (s *Server) render(name string, data interface{}) ([]byte, error) {
	t, err := template.ParseFiles(filepath.Join(s.templates, name+".html"))
	if err != nil {
		return nil, err
	}
	var br bytes.Buffer
	if err = t.Execute(&br, data); err != nil {
		return nil, err
	}
	return br.Bytes(), nil
}
