package applications

import (
	"context"
	"errors"

	"github.com/gabrielmatsan/teste-api/cmd/db"
	"github.com/gabrielmatsan/teste-api/internal/helpers"
	"github.com/gabrielmatsan/teste-api/internal/shared/email/templates"
	"github.com/gabrielmatsan/teste-api/internal/shared/singlaton"
	"github.com/gabrielmatsan/teste-api/internal/user/model"
)

func CreateUserUseCase(ctx context.Context, dataUser *model.CreateUser) (model.User, error) {

	// repo instance
	userRepository := singlaton.GetUserRepository()

	isEmailAlreadyUsed, err := userRepository.GetUserByEmail(ctx, dataUser.Email)

	if err != nil {
		return model.User{}, err
	}

	if isEmailAlreadyUsed != nil {
		return model.User{}, errors.New("email already used")
	}

	// begin transaction
	db := db.GetDB()
	tx, err := helpers.BeginTransaction(ctx, db)

	if err != nil {
		return model.User{}, err
	}

	// defer rollback if error
	var finalErr error
	defer func() {
		if finalErr != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// create user
	user, err := userRepository.CreateUser(ctx, dataUser, tx)

	if err != nil {
		finalErr = err
		return model.User{}, err
	}

	// send to sqs + lambda
	go templates.SendWelcomeEmail(user)

	return user, nil
}
