package main

import (
	"fitnes-lessons/internal/app"
	"fitnes-lessons/internal/config"
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
	application := app.New(logger, &appConfig)
	go func() {
		application.GRPCServer.MustRun()
	}()
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
	return log
}
