package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type User struct {
	ID        string     `json:"id" db:"id" validate:"required"`
	FirstName string     `json:"firstName" db:"first_name" validate:"required"`
	LastName  string     `json:"lastName" db:"last_name" validate:"required"`
	Email     string     `json:"email" db:"email" validate:"required,email"`
	Password  string     `json:"password" db:"password" validate:"required"`
	Role      Role       `json:"role" db:"role" validate:"required"`
	Status    Status     `json:"status" db:"status" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at" validate:"required"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at" validate:"omitempty"`
	IsActive  bool       `json:"isActive" db:"is_active" validate:"required"`
}

type CreateUser struct {
	FirstName string `json:"firstName" db:"first_name" validate:"required"`
	LastName  string `json:"lastName" db:"last_name" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required,email"`
	Password  string `json:"password" db:"password" validate:"required"`
}
