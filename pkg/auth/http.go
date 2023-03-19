package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/matryer/way"
	"github.com/pkg/errors"
)

func NewServer(service Service, logger log.Logger) http.Handler {
	s := server{service: service}

	var handleLogin http.Handler
	handleLogin = s.handleLogin()
	handleLogin = httpLoggingMiddleware(logger, "handleLogin")(handleLogin)

	var handleLogout http.Handler
	handleLogout = s.handleLogout()
	handleLogout = httpLoggingMiddleware(logger, "handleLogout")(handleLogout)

	router := way.NewRouter()

	router.Handle("POST", "/auth/v1/login", handleLogin)
	router.Handle("POST", "/auth/v1/logout", handleLogout)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, ErrResourceNotFound) })

	return router

}

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

type ErrInvalidRequestBody struct{ err error }

func (e ErrInvalidRequestBody) Error() string { return fmt.Sprintf("invalid request body: %v", e.err) }

type server struct {
	service Service
}

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		EmployeeID string `json:"employeeid"`
		Password   string `json:"password"`
	}
	type response struct {
		ID          int64     `json:"id"`
		Name        string    `json:"name"`
		Designation string    `json:"designation"`
		EmployeeID  string    `json:"employeeid"`
		CreatedAt   time.Time `json:"createdAt"`
		Token       string    `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		lr, err := s.service.Login(r.Context(), LoginRequest{EmployeeID: req.EmployeeID, Password: req.Password})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set(contentTypeKey, contentTypeValue)
		json.NewEncoder(w).Encode(response{ID: lr.User.ID, Name: lr.User.Name, Designation: lr.User.Designation, EmployeeID: lr.User.EmployeeID, CreatedAt: lr.User.CreatedAt, Token: lr.Token})
	}
}

func (s *server) handleLogout() http.HandlerFunc {
	type request struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, ErrInvalidRequestBody{err})
			return
		}

		err := s.service.Logout(r.Context(), req.Token)

		if err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set(contentTypeKey, contentTypeValue)

	switch err {
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
