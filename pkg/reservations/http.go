package reservations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/markhaur/messapp-backend/pkg"
	"github.com/matryer/way"
)

func NewServer(service Service, logger log.Logger) http.Handler {
	s := server{service: service}

	var handleSaveReservation http.Handler
	handleSaveReservation = s.handleSaveReservation()
	handleSaveReservation = httpLoggingMiddleware(logger, "handleSaveReservation")(handleSaveReservation)

	var handleListReservations http.Handler
	handleListReservations = s.handleListReservations()
	handleListReservations = httpLoggingMiddleware(logger, "handleListReservations")(handleListReservations)

	var handleGetReservations http.Handler
	handleGetReservations = s.handleGetReservation()
	handleGetReservations = httpLoggingMiddleware(logger, "handleGetReservations")(handleGetReservations)

	var handleGetReservationsByID http.Handler
	handleGetReservationsByID = s.handleGetReservationByID()
	handleGetReservationsByID = httpLoggingMiddleware(logger, "handleGetReservations")(handleGetReservationsByID)

	var handleRemoveReservation http.Handler
	handleRemoveReservation = s.handleRemoveReservation()
	handleRemoveReservation = httpLoggingMiddleware(logger, "handleRemoveReservation")(handleRemoveReservation)

	var handleUpdateReservation http.Handler
	handleUpdateReservation = s.handleUpdateReservation()
	handleUpdateReservation = httpLoggingMiddleware(logger, "handleUpdateReservation")(handleUpdateReservation)

	router := way.NewRouter()

	router.Handle("POST", "/resvlist/v1/reservations", handleSaveReservation)
	router.Handle("GET", "/resvlist/v1/reservations", handleListReservations)
	router.Handle("GET", "/resvlist/v1/reservations/:date", handleGetReservations)
	// router.Handle("GET", "/resvlist/v1/reservations/:id", handleGetReservations)
	router.Handle("GET", "/resvlist/v1/reservationsbyid/:user_id", handleGetReservationsByID)
	router.Handle("DELETE", "/resvlist/v1/reservation/:id", handleRemoveReservation)
	router.Handle("PUT", "/resvlist/v1/reservation/:id", handleUpdateReservation)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, ErrResourceNotFound) })

	return router
}

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"
)

var (
	ErrNonNumericReservationID = errors.New("reservation id in path must be numberic")
	ErrResourceNotFound        = errors.New("resource not found")
	ErrMethodNotAllowed        = errors.New("method not allowed")
)

type ErrInvalidRequestBody struct{ err error }

func (e ErrInvalidRequestBody) Error() string { return fmt.Sprintf("invalid request body: %v", e.err) }

type server struct {
	service Service
}

func (s *server) handleSaveReservation() http.HandlerFunc {
	type request struct {
		UserID          int64               `json:"user_id"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
	}
	type response struct {
		ID              int64               `json:"id"`
		UserID          int64               `json:"user_id"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
		CreatedAt       time.Time           `json:"createdAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		reservation, err := s.service.Save(r.Context(), pkg.Reservation{UserID: req.UserID, ReservationTime: req.ReservationTime, Type: req.Type, NoOfGuests: req.NoOfGuests})
		if err != nil {
			writeError(w, err)
			return
		}
		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(response{ID: reservation.ID, UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: reservation.Type, NoOfGuests: req.NoOfGuests, CreatedAt: reservation.CreatedAt})
	}
}

func (s *server) handleGetReservation() http.HandlerFunc {
	type reservation struct {
		ID              int64               `json:"id"`
		UserID          int64               `json:"user_id"`
		Name            string              `json:"name"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
		CreatedAt       time.Time           `json:"createdAt"`
	}
	type response []reservation
	return func(w http.ResponseWriter, r *http.Request) {
		var reservations []pkg.Reservation
		// id, err := strconv.ParseInt(way.Param(r.Context(), "id"), 10, 64)
		// if err == nil {
		// 	resv, err := s.service.FindByID(r.Context(), id)
		// 	if err != nil {
		// 		writeError(w, err)
		// 		return
		// 	}
		// 	reservations = append(reservations, *resv)
		// }

		user_id, err := strconv.ParseInt(way.Param(r.Context(), "user_id"), 10, 64)
		if err == nil {
			resvs, err := s.service.FindByEmployeeID(r.Context(), user_id)
			if err != nil {
				writeError(w, err)
				return
			}
			reservations = resvs
		}

		format := "2006-01-02"
		date, err := time.Parse(format, way.Param(r.Context(), "date"))
		if err == nil {
			resvs, err := s.service.FindByDate(r.Context(), date)
			if err != nil {
				writeError(w, err)
				return
			}
			reservations = resvs
		}

		resp := make(response, 0, len(reservations))

		for _, resv := range reservations {
			resp = append(resp, reservation{ID: resv.ID, UserID: resv.UserID, Name: resv.Name, ReservationTime: resv.ReservationTime, Type: resv.Type, NoOfGuests: resv.NoOfGuests, CreatedAt: resv.CreatedAt})
		}

		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *server) handleGetReservationByID() http.HandlerFunc {
	type reservation struct {
		ID              int64               `json:"id"`
		UserID          int64               `json:"user_id"`
		Name            string              `json:"name"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
		CreatedAt       time.Time           `json:"createdAt"`
	}
	type response []reservation
	return func(w http.ResponseWriter, r *http.Request) {
		var reservations []pkg.Reservation

		user_id, err := strconv.ParseInt(way.Param(r.Context(), "user_id"), 10, 64)
		if err == nil {
			resvs, err := s.service.FindByEmployeeID(r.Context(), user_id)
			if err != nil {
				writeError(w, err)
				return
			}
			reservations = resvs
		}

		resp := make(response, 0, len(reservations))
		for _, resv := range reservations {
			resp = append(resp, reservation{ID: resv.ID, UserID: resv.UserID, Name: resv.Name, ReservationTime: resv.ReservationTime, Type: resv.Type, NoOfGuests: resv.NoOfGuests, CreatedAt: resv.CreatedAt})
		}

		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *server) handleListReservations() http.HandlerFunc {
	type reservation struct {
		ID              int64               `json:"id"`
		UserID          int64               `json:"user_id"`
		Name            string              `json:"name"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
		CreatedAt       time.Time           `json:"createdAt"`
	}
	type response []reservation

	return func(w http.ResponseWriter, r *http.Request) {
		list, err := s.service.List(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}

		resp := make(response, 0, len(list))
		for _, v := range list {
			resp = append(resp, reservation{ID: v.ID, UserID: v.UserID, Name: v.Name, ReservationTime: v.ReservationTime, Type: v.Type, NoOfGuests: v.NoOfGuests, CreatedAt: v.CreatedAt})
		}
		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *server) handleRemoveReservation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(way.Param(r.Context(), "id"), 10, 64)
		if err != nil {
			writeError(w, ErrNonNumericReservationID)
			return
		}

		if err := s.service.Remove(r.Context(), id); err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleUpdateReservation() http.HandlerFunc {
	type request struct {
		UserID          int64               `json:"user_id"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
	}
	type response struct {
		ID              int64               `json:"id"`
		UserID          int64               `json:"user_id"`
		ReservationTime time.Time           `json:"reservation_time"`
		Type            pkg.ReservationType `json:"type"`
		NoOfGuests      int64               `json:"no_of_guests"`
		CreatedAt       time.Time           `json:"createdAt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(way.Param(r.Context(), "id"), 10, 64)
		if err != nil {
			writeError(w, ErrNonNumericReservationID)
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		reservation, isCreated, err := s.service.Update(r.Context(), pkg.Reservation{ID: id, UserID: req.UserID, ReservationTime: req.ReservationTime, Type: req.Type, NoOfGuests: req.NoOfGuests})
		if err != nil {
			writeError(w, err)
			return
		}

		if isCreated {
			w.WriteHeader(http.StatusCreated)
		}

		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(response{ID: reservation.ID, UserID: reservation.UserID, ReservationTime: reservation.ReservationTime, Type: reservation.Type, NoOfGuests: reservation.NoOfGuests, CreatedAt: reservation.CreatedAt})
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set(contentTypeKey, contentTypeValue)

	switch err {
	case ErrResourceNotFound, pkg.ErrReservationNotFound:
		w.WriteHeader(http.StatusNotFound)
	case pkg.ErrReservationAlreadyExists:
		w.WriteHeader(http.StatusConflict)
	case ErrNonNumericReservationID:
		w.WriteHeader(http.StatusBadRequest)
	case ErrMethodNotAllowed:
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		switch err.(type) {
		case ErrInvalidRequestBody:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func httpLoggingMiddleware(logger log.Logger, operation string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()
			lrw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lrw, r)
			logger.Log(
				"operation", operation,
				"method", r.Method,
				"path", r.URL.Path,
				"took", time.Since(begin),
				"status", lrw.statusCode,
			)
		})
	}
}
