package userlist

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/markhaur/messapp-backend/pkg"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(s Service) Service { return &loggingMiddleware{logger, s} }
}

type loggingMiddleware struct {
	logger log.Logger
	Service
}

func (s *loggingMiddleware) Save(ctx context.Context, user pkg.User) (_ *pkg.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "save",
			"name", user.Name,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.Save(ctx, user)
}

func (s *loggingMiddleware) List(ctx context.Context) (_ []pkg.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "list",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.List(ctx)
}

func (s *loggingMiddleware) Remove(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "remove",
			"id", id,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.Remove(ctx, id)
}

func (s *loggingMiddleware) Update(ctx context.Context, user pkg.User) (_ *pkg.User, _ bool, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "update",
			"name", user.Name,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.Update(ctx, user)
}
