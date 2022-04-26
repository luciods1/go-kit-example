package services

import (
	"context"

	"github.com/go-kit/log"
)

//	errors definition

type Service interface {
	HealthCheck(ctx context.Context) (HealthCheckResponse, error)
}

type healthService struct {
	//here will be all services and connection dependencies dependencies
}

func New(logger log.Logger) Service {
	var svc Service
	svc = NewHealthService()
	svc = LoggingMiddleware(logger)(svc)
	return svc
}

func (h healthService) HealthCheck(ctx context.Context) (HealthCheckResponse, error) {
	return HealthCheckResponse{
		Connected: true,
		Name:      "Service Name",
	}, nil
}

func NewHealthService() Service {
	return healthService{}
}

type HealthCheckResponse struct {
	Connected bool
	Name      string
	Err       error `json:"-"`
}
