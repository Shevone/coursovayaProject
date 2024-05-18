package main

import (
	"fitnes-account/internal/app"
	"fitnes-account/internal/config"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func envLoad() {
	// Инициализируем перменные окржения, с помощью .env файла
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is empty")
	}
	appSalt := os.Getenv("APP_SECRET")
	if appSalt == "" {
		panic("salt for jwt is empty")
	}
}
func main() {
	envLoad()

	cfg := config.MustLoad()
	log := setupLogger()

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
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
