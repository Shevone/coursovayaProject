package main

import (
	"fitnes-account/internal/app"
	"fitnes-account/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//envLoad()

	appConfig := app.Config{}
	config.MustConfig(&appConfig)

	log := setupLogger()
	slog.SetDefault(log)

	application := app.New(log, &appConfig)
	go func() {
		application.GRPCServer.MustRun()
	}()
	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")

}
func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
