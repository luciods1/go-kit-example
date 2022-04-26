package services

import (
	"context"

	"github.com/go-kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

func (mw loggingMiddleware) HealthCheck(ctx context.Context) (HealthCheckResponse, error) {
	// log interestings things here
	return mw.next.HealthCheck(ctx)
}
