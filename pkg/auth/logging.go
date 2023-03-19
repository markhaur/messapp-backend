package auth

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

func LoggingMiddleware(logger log.Logger) Middlewware {
	return func(s Service) Service { return &loggingMiddleware{logger, s} }
}

type loggingMiddleware struct {
	logger log.Logger
	Service
}

func (l *loggingMiddleware) Login(ctx context.Context, request LoginRequest) (_ *LoginResponse, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "login",
			"employee_id", request.EmployeeID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Login(ctx, request)
}

func (l *loggingMiddleware) Logout(ctx context.Context, token string) (err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "logout",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Logout(ctx, token)
}
