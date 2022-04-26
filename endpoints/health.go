package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/luciodesimone/go-kit-example/services"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	HealthEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(hSvc services.Service) Set {
	return Set{
		HealthEndpoint: MakeHealthEndpoints(hSvc),
	}
}

func MakeHealthEndpoints(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, err := svc.HealthCheck(ctx)

		if err != nil {
			return r, err
		}

		return r, nil
	}
}
func New(svc services.Service, logger log.Logger, duration metrics.Histogram) Set {
	var healthEndpoint endpoint.Endpoint

	healthEndpoint = MakeHealthEndpoints(svc)
	healthEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(healthEndpoint)
	healthEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(healthEndpoint)
	healthEndpoint = LoggingMiddleware(log.With(logger, "method", "Health"))(healthEndpoint)

	// TODO: set up instrumentation middleware in case you include some one
	// healthEndpoint = InstrumentingMiddleware(duration.With("method", "Health"))(healthEndpoint)

	return Set{
		HealthEndpoint: healthEndpoint,
	}
}

type HealthCheckRequest struct{}
