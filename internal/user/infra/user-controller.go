package infra

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabrielmatsan/teste-api/internal/user/applications"
	"github.com/gabrielmatsan/teste-api/internal/user/model"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserModel model.CreateUser

	if err := json.NewDecoder(r.Body).Decode(&createUserModel); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	user, err := applications.CreateUserUseCase(r.Context(), &createUserModel)
	if err != nil {
		// Verifica se é erro de email já usado
		if err.Error() == "email already used" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUserController() *UserController {
	return &UserController{}
}
