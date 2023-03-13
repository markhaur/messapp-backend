package auth

import (
	"context"

	"github.com/markhaur/messapp-backend/pkg"
)

type LoginResponse struct {
	User  pkg.User
	Token []byte
}

type LoginRequest struct {
	Username string
	Password []byte
}

type Service interface {
	Login(context.Context, LoginRequest) (LoginResponse, error)
	Logout(ctx context.Context, token string) error
}
