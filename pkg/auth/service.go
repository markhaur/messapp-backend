package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markhaur/messapp-backend/pkg"
	"github.com/pkg/errors"
)

var tokenKey = []byte("secret_key")

type claim struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Designation string    `json:"designation"`
	EmployeeID  string    `json:"employeeid"`
	CreatedAt   time.Time `json:"createdAt"`
	jwt.StandardClaims
}

type LoginResponse struct {
	User  pkg.User
	Token string
}

type LoginRequest struct {
	EmployeeID string
	Password   string
}

type Service interface {
	Login(context.Context, LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
}

type Middlewware func(Service) Service

type service struct {
	repository pkg.UserRepository
}

func NewService(repo pkg.UserRepository) Service {
	return &service{repository: repo}
}

func (s *service) Login(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	user, err := s.repository.FindByEmployeeID(ctx, request.EmployeeID)
	if err != nil {
		return nil, err
	}

	if user.Password != request.Password {
		return nil, errors.New("invalid password")
	}

	token, err := createToken(*user)
	if err != nil {
		return nil, err
	}

	response := LoginResponse{
		User:  *user,
		Token: token,
	}
	return &response, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	// TODO: need to implement it
	return nil
}

func createToken(user pkg.User) (string, error) {
	// set the expiration time
	expirationTime := time.Now().Add(1 * time.Hour)

	// creating jwt claim
	claim := &claim{
		ID:          user.ID,
		Name:        user.Name,
		Designation: user.Designation,
		EmployeeID:  user.EmployeeID,
		CreatedAt:   user.CreatedAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create jwt token using claims and signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(tokenKey)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
