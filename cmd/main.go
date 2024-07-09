package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"time-tracker/internal/config"
	usersHandler "time-tracker/internal/controller/user"
	"time-tracker/internal/lib/logger"
	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/repository/externalapi"
	storage "time-tracker/internal/repository/postgres"
	taskService "time-tracker/internal/service/task"
	usersService "time-tracker/internal/service/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Info("initializing server...", slog.String("port", cfg.Server.Port))

	// Data layer
	storage, err := storage.New(cfg.Storage)
	if err != nil {
		log.Error("storage initial error", sl.Error(err))
		return
	}

	externalAPI := externalapi.New(cfg.ExternalAPI)

	// Service layer
	usersService := usersService.New(storage, externalAPI, log)
	taskService := taskService.New(storage, log)

	// Controllers layer
	usersHandler := usersHandler.New(usersService, taskService, log)

	// Init router
	r := chi.NewRouter()
	chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", usersHandler.Register())

	// Init server
	srv := http.Server{
		Handler:      r,
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.Timeout * time.Second,
		WriteTimeout: cfg.Server.Timeout * time.Second,
		IdleTimeout:  cfg.Server.Timeout * time.Second,
	}

	log.Info("server initialized")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("server error", sl.Error(err))
		}
	}()

	log.Info("server is running...")

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", sl.Error(err))
	}

	storage.Close()

	log.Info("server stopped")
}
