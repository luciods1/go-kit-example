package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/luciodesimone/go-kit-example/endpoints"
	"github.com/luciodesimone/go-kit-example/services"
	"github.com/luciodesimone/go-kit-example/transports"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	if err := run(logger); err != nil {
		logger.Log("Something went wrong")
		os.Exit(1)
	}
}

func run(logger log.Logger) error {
	var duration metrics.Histogram

	// Example of metrics using Prometheus client (assuming counters where already created)
	// duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	// 	Namespace: "example",
	// 	Subsystem: "healthsvc",
	// 	Name:      "request_duration_seconds",
	// 	Help:      "Request duration in seconds.",
	// }, []string{"method", "success"})

	service := services.New(logger)
	endpoints := endpoints.New(service, logger, duration)
	httpHandler := transports.NewHTTPHandler(endpoints, logger)

	port := os.Getenv("HTTP_PORT")
	srv := http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      httpHandler,
	}
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		return err
	case shutdownSignal := <-shutdown:
		fmt.Println("Starting graceful shutdown...")

		bkg := context.Background()
		ctx, cancel := context.WithTimeout(bkg, time.Second*30)
		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			err = srv.Close()
		}

		switch {
		case shutdownSignal == syscall.SIGSTOP:
			return errors.New("Process has been canceled (SIGSTOP)")
		case err != nil:
			return errors.New("Could not stop server gracefully")
		}
	}

	return nil
}
