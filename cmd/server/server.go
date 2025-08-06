package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gabrielmatsan/teste-api/internal/user/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
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
	s.router.Route("/api/v1", func(r chi.Router) {
		routes.UserRoutes(r)
	})
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL completa
	))

	chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Rota registrada: %s %s", method, route)
		return nil
	})

}

// Start inicia o servidor HTTP.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Servidor Chi iniciado na porta %s", s.port)
	return http.ListenAndServe(addr, s.router)
}
