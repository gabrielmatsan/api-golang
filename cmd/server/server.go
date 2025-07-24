package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	router *chi.Mux
	db     *sqlx.DB
	port   string
}

func NewServer(db *sqlx.DB, port string) *Server {
	r := chi.NewRouter()

	return &Server{
		router: r,
		db:     db,
		port:   port,
	}
}

// Setup para middlewares (esses considerados globais, autenticação sendo feita em cada rota)
func (s *Server) SetupMiddlewares() {
	s.router.Use(middleware.RequestID) // Gera um ID único para cada requisição
	s.router.Use(middleware.Logger)    // Loga as requisições
	s.router.Use(middleware.Recoverer) // Recupera de panics e loga o stack trace
}

func (s *Server) SetupRoutes() {
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
}

// Start inicia o servidor HTTP.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Servidor Chi iniciado na porta %s", s.port)
	return http.ListenAndServe(addr, s.router)
}
