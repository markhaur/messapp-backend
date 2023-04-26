package userlist

import (
	"context"
	"fmt"
	"time"

	"github.com/markhaur/messapp-backend/pkg"
)

type Service interface {
	Save(context.Context, pkg.User) (*pkg.User, error)
	List(context.Context) ([]pkg.User, error)
	Update(context.Context, pkg.User) (*pkg.User, bool, error)
	Remove(context.Context, int64) error
}

// Middleware describes a Service Middleware
type Middleware func(Service) Service

type service struct {
	repository pkg.UserRepository
}

func NewService(repository pkg.UserRepository) Service {
	return &service{repository: repository}
}

func (s *service) Save(ctx context.Context, user pkg.User) (*pkg.User, error) {
	user.CreatedAt = time.Now()
	user.Password = "password@1234"
	if err := s.repository.Insert(ctx, &user); err != nil {
		return nil, fmt.Errorf("could not save user: %v", err)
	}
	return &user, nil
}

func (s service) List(ctx context.Context) ([]pkg.User, error) {
	list, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list all users: %v", err)
	}
	return list, nil
}

func (s *service) Update(ctx context.Context, user pkg.User) (*pkg.User, bool, error) {
	err := s.repository.Update(ctx, &user)
	if err == pkg.ErrUserNotFound {
		err = s.repository.Insert(ctx, &user)
		if err != nil {
			return nil, false, fmt.Errorf("could not create user: %v", err)
		}
		return &user, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("could not update user: %v", err)
	}
	return &user, false, nil
}

func (s *service) Remove(ctx context.Context, id int64) error {
	if err := s.repository.DeleteByID(ctx, id); err != nil {
		if err == pkg.ErrUserNotFound {
			return err
		}
		return fmt.Errorf("could not remove user: %v", err)
	}
	return nil
}
