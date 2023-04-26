package pkg

import (
	"context"
	"errors"
	"time"
)

var (
	ErrReservationNotFound      = errors.New("reservation not found")
	ErrReservationAlreadyExists = errors.New("reservation already exists")
)

type ReservationType int

const (
	Breakfast ReservationType = iota + 1
	Lunch
	Dinner
)

func (r ReservationType) String() string {
	types := [...]string{"breakfast", "lunch", "dinner"}
	if r < Breakfast || r > Dinner {
		return ""
	}
	return types[r-1]
}

type Reservation struct {
	ID              int64
	UserID          int64
	Name            string
	ReservationTime time.Time
	Type            ReservationType
	NoOfGuests      int64
	CreatedAt       time.Time
}

type ReservationRepository interface {
	Insert(context.Context, *Reservation) error
	FindAll(context.Context) ([]Reservation, error)
	FindByID(context.Context, int64) (*Reservation, error)
	FindByEmployeeID(context.Context, int64) ([]Reservation, error)
	FindByDate(context.Context, time.Time) ([]Reservation, error)
	Update(context.Context, *Reservation) error
	DeleteByID(context.Context, int64) error
}
