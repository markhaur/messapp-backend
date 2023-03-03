package userlist

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

	var handleSaveUser http.Handler
	handleSaveUser = s.handleSaveUser()
	handleSaveUser = httpLoggingMiddleware(logger, "handleSaveUser")(handleSaveUser)

	var handleListUsers http.Handler
	handleListUsers = s.handleListUsers()
	handleListUsers = httpLoggingMiddleware(logger, "handleListUsers")(handleListUsers)

	var handleRemoveUser http.Handler
	handleRemoveUser = s.handleRemoveUser()
	handleListUsers = httpLoggingMiddleware(logger, "handleRemoveUser")(handleRemoveUser)

	var handleUpdateUser http.Handler
	handleUpdateUser = s.handleUpdateUser()
	handleUpdateUser = httpLoggingMiddleware(logger, "handleUpdateUser")(handleUpdateUser)

	router := way.NewRouter()

	router.Handle("POST", "/v1/users", handleSaveUser)
	router.Handle("GET", "/v1/users", handleListUsers)
	router.Handle("DELETE", "/v1/user/:id", handleRemoveUser)
	router.Handle("PUT", "/v1/user/:id", handleUpdateUser)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, ErrResourceNotFound) })

	return router
}

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"
)

var (
	ErrNonNumericUserID = errors.New("user id in path must be numberic")
	ErrResourceNotFound = errors.New("resource not found")
	ErrMethodNotAllowed = errors.New("method not allowed")
)

type ErrInvalidRequestBody struct{ err error }

func (e ErrInvalidRequestBody) Error() string { return fmt.Sprintf("invalid request body: %v", e.err) }

type server struct {
	service Service
}

func (s *server) handleSaveUser() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Password    string `json:"password"`
		Designation string `json:"designation"`
		EmployeeID  string `json:"employeeid"`
	}
	type response struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Designation string    `json:"designation"`
		EmployeeID  string    `json:"employeeID"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		user, err := s.service.Save(r.Context(), pkg.User{Name: req.Name, Designation: req.Designation, EmployeeID: req.EmployeeID})
		if err != nil {
			writeError(w, err)
			return
		}
		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(response{ID: user.ID, Name: user.Name, Designation: user.Designation, EmployeeID: user.EmployeeID, CreatedAt: user.CreatedAt})
	}
}

func (s *server) handleListUsers() http.HandlerFunc {
	type user struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Designation string    `json:"designation"`
		EmployeeID  string    `json:"employeeID"`
		CreatedAt   time.Time `json:"createdAt"`
	}
	type response []user

	return func(w http.ResponseWriter, r *http.Request) {
		list, err := s.service.List(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}

		resp := make(response, 0, len(list))
		for _, v := range list {
			resp = append(resp, user{ID: v.ID, Name: v.Name, Designation: v.Designation, EmployeeID: v.EmployeeID, CreatedAt: v.CreatedAt})
		}
		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *server) handleRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(way.Param(r.Context(), "id"), 10, 64)
		if err != nil {
			writeError(w, ErrNonNumericUserID)
			return
		}

		if err := s.service.Remove(r.Context(), id); err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) handleUpdateUser() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		Password    string `json:"password"`
		Designation string `json:"designation"`
		EmployeeID  string `json:"employeeid"`
	}
	type response struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Designation string    `json:"designation"`
		EmployeeID  string    `json:"employeeID"`
		CreatedAt   time.Time `json:"createdAt"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(way.Param(r.Context(), "id"), 10, 64)
		if err != nil {
			writeError(w, ErrNonNumericUserID)
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		user, isCreated, err := s.service.Update(r.Context(), pkg.User{ID: id, Name: req.Name, Password: req.Password, Designation: req.Designation, EmployeeID: req.EmployeeID})
		if err != nil {
			writeError(w, err)
			return
		}

		if isCreated {
			w.WriteHeader(http.StatusCreated)
		}

		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(response{ID: user.ID, Name: user.Name, Designation: user.Designation, EmployeeID: user.EmployeeID})
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set(contentTypeKey, contentTypeValue)

	switch err {
	case ErrResourceNotFound, pkg.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	case pkg.ErrUserAlreadyExists:
		w.WriteHeader(http.StatusConflict)
	case ErrNonNumericUserID:
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
