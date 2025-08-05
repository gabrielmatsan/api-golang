package applications

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabrielmatsan/teste-api/internal/shared/email/templates"
	"github.com/gabrielmatsan/teste-api/internal/shared/singlaton"
	"github.com/gabrielmatsan/teste-api/internal/user/model"
)

func CreateUserUseCase(ctx context.Context, dataUser *model.CreateUser) (model.User, error) {

	fmt.Println("CreateUserUseCase")
	userRepository := singlaton.GetUserRepository()

	fmt.Println("dataUser", dataUser)

	isEmailAlreadyUsed, err := userRepository.GetUserByEmail(ctx, dataUser.Email)

	if err != nil {
		return model.User{}, err
	}

	if isEmailAlreadyUsed != nil {
		return model.User{}, errors.New("email already used")
	}

	user, err := userRepository.CreateUser(ctx, dataUser)

	if err != nil {
		return model.User{}, err
	}

	go templates.SendWelcomeEmail(user)

	return user, nil
}
