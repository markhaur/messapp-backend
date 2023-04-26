package pkg

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID          int64
	Name        string
	Password    string
	Designation string
	EmployeeID  string
	IsAdmin     int32
	IsActive    int32
	CreatedAt   time.Time
}

type UserRepository interface {
	Insert(context.Context, *User) error
	FindAll(context.Context) ([]User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByEmployeeID(context.Context, string) (*User, error)
	Update(context.Context, *User) error
	DeleteByID(context.Context, int64) error
}
