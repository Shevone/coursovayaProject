package main

import (
	"context"
	"fitnes-lessons/internal/app"
	"fitnes-lessons/internal/config"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	appConfig := app.Config{}
	config.MustConfig(&appConfig)

	logger := setupLogger()
	slog.SetDefault(logger)

	errWg, errCtx := errgroup.WithContext(context.Background())
	application := app.New(errCtx, logger, &appConfig)
	errWg.Go(func() error {
		return application.GRPCServer.Run()
	})
	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("Gracefully stopped")
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(log)
	return log
}
