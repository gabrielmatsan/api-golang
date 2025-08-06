package routes

import (
	"github.com/gabrielmatsan/teste-api/internal/user/infra"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router) {
	userController := infra.GetUserController()


	
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.CreateUser)
	})
}
