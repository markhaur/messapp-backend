package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/markhaur/messapp-backend/pkg"
	"github.com/markhaur/messapp-backend/pkg/mysql/gen"
)

type reservationRepository struct {
	queries *gen.Queries
}

func NewReservationRepository(db *sql.DB) pkg.ReservationRepository {
	return &reservationRepository{queries: gen.New(db)}
}

func (r *reservationRepository) Insert(ctx context.Context, reservation *pkg.Reservation) error {
	inserted, err := r.queries.CreateReservation(ctx, gen.CreateReservationParams{UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: int64(reservation.Type), NoOfGuests: reservation.NoOfGuests, CreatedAt: time.Now()})
	if err != nil {
		return err
	}
	reservation.ID, _ = inserted.LastInsertId()
	reservation.CreatedAt = reservation.CreatedAt.In(time.Local)
	return nil
}

func (r *reservationRepository) FindAll(context.Context) ([]pkg.Reservation, error) {
	return nil, nil
}

func (r *reservationRepository) FindByID(ctx context.Context, id int64) (*pkg.Reservation, error) {
	reservation, err := r.queries.GetReservationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pkg.Reservation{ID: reservation.ID, UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: pkg.ReservationType(reservation.Type), NoOfGuests: reservation.NoOfGuests, CreatedAt: reservation.CreatedAt}, nil
}

func (r *reservationRepository) FindByEmployeeID(ctx context.Context, employee_id int64) ([]pkg.Reservation, error) {
	reservations, err := r.queries.GetReservationsByEmployeeID(ctx, employee_id)
	if err != nil {
		return nil, err
	}

	var list []pkg.Reservation
	for _, reservation := range reservations {
		list = append(list, pkg.Reservation{ID: reservation.ID, UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: pkg.ReservationType(reservation.Type), NoOfGuests: reservation.NoOfGuests, CreatedAt: reservation.CreatedAt})
	}
	return list, nil
}

func (r *reservationRepository) FindByDate(ctx context.Context, date time.Time) ([]pkg.Reservation, error) {
	reservations, err := r.queries.GetReservationsByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	var list []pkg.Reservation
	for _, reservation := range reservations {
		list = append(list, pkg.Reservation{ID: reservation.ID, UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: pkg.ReservationType(reservation.Type), NoOfGuests: reservation.NoOfGuests, CreatedAt: reservation.CreatedAt})
	}
	return list, nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *pkg.Reservation) error {
	return r.queries.UpdateReservation(ctx, gen.UpdateReservationParams{ID: reservation.ID, NoOfGuests: reservation.NoOfGuests})
}

func (r *reservationRepository) DeleteByID(ctx context.Context, id int64) error {
	return r.queries.DeleteReservation(ctx, id)
}
