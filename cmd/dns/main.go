package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"space2218/internal/dns"
	httpapi "space2218/internal/dns/ports/http"
	"space2218/internal/dns/services"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("started")

	cfg := LoadConfig()
	logger.Info("config loaded", zap.Any("config", cfg))

	// system signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// services
	navigationService := services.NewNavigation(dns.SectorID(cfg.SectorID))

	// server
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: httpapi.Routes(logger, navigationService),
	}

	logger.Info("starting server", zap.String("address", cfg.Address))
	go func() {
		err := server.ListenAndServe()
		switch {
		case err == nil:
		// impossible branch
		case errors.Is(err, http.ErrServerClosed):
			logger.Info("server closed")
		default:
			logger.Error("server.ListenAndServe", zap.Error(err))
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server.Shutdown", zap.Error(err))
	}

	logger.Info("finished")
}
