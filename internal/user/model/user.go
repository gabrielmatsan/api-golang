package model

import "database/sql"

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
