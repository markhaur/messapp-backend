package reservations

import (
	"context"
	"fmt"
	"time"

	"github.com/markhaur/messapp-backend/pkg"
)

type Service interface {
	Save(context.Context, pkg.Reservation) (*pkg.Reservation, error)
	List(context.Context) ([]pkg.Reservation, error)
	FindByID(context.Context, int64) (*pkg.Reservation, error)
	FindByEmployeeID(context.Context, int64) ([]pkg.Reservation, error)
	FindByDate(context.Context, time.Time) ([]pkg.Reservation, error)
	// TODO: need to remove
	Update(context.Context, pkg.Reservation) (*pkg.Reservation, bool, error)
	Remove(context.Context, int64) error
}

// Middleware describes a Service Middleware
type Middleware func(Service) Service

type service struct {
	repository pkg.ReservationRepository
}

func NewService(repository pkg.ReservationRepository) Service {
	return &service{repository: repository}
}

func (s *service) Save(ctx context.Context, reservation pkg.Reservation) (*pkg.Reservation, error) {
	if err := s.repository.Insert(ctx, &reservation); err != nil {
		return nil, fmt.Errorf("could not save reservation: %v", err)
	}
	return &reservation, nil
}

func (s *service) FindByID(ctx context.Context, id int64) (*pkg.Reservation, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByEmployeeID(ctx context.Context, id int64) ([]pkg.Reservation, error) {
	return s.repository.FindByEmployeeID(ctx, id)
}

func (s *service) FindByDate(ctx context.Context, date time.Time) ([]pkg.Reservation, error) {
	return s.repository.FindByDate(ctx, date)
}

func (s service) List(ctx context.Context) ([]pkg.Reservation, error) {
	list, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list all reservations: %v", err)
	}
	return list, nil
}

func (s *service) Update(ctx context.Context, reservation pkg.Reservation) (*pkg.Reservation, bool, error) {
	err := s.repository.Update(ctx, &reservation)
	if err == pkg.ErrReservationNotFound {
		err = s.repository.Insert(ctx, &reservation)
		if err != nil {
			return nil, false, fmt.Errorf("could not create reservation: %v", err)
		}
		return &reservation, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("could not update reservation: %v", err)
	}
	return &reservation, false, nil
}

func (s *service) Remove(ctx context.Context, id int64) error {
	if err := s.repository.DeleteByID(ctx, id); err != nil {
		if err == pkg.ErrReservationNotFound {
			return err
		}
		return fmt.Errorf("could not remove reservation: %v", err)
	}
	return nil
}
