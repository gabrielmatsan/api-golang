package repository

import (
	"context"
	"database/sql"

	"github.com/gabrielmatsan/teste-api/internal/user/model"
	"github.com/jmoiron/sqlx"
	"github.com/nrednav/cuid2"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, dataUser *model.CreateUser, tx *sqlx.Tx) (model.User, error) {

	id := cuid2.Generate()

	query := `
		INSERT INTO users (id, first_name, last_name, email, password, role, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *
	`

	args := []any{id, dataUser.FirstName, dataUser.LastName, dataUser.Email, dataUser.Password, "user", "active"}

	var user model.User

	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &user, query, args...)
	} else {
		err = r.db.GetContext(ctx, &user, query, args...)
	}

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {

	query := `
		SELECT *
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
