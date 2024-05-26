package main

import (
	"context"
	"fitnes-gateway/internal/config"
	handler "fitnes-gateway/internal/handler"
	"fitnes-gateway/internal/lib"
	"fitnes-gateway/internal/server"
	service "fitnes-gateway/internal/service"
	"fitnes-gateway/internal/service/clients/accounts"
	"fitnes-gateway/internal/service/clients/lessons"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	Clients    ClientsConfig `yaml:"clients"`
	AppSecret  string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type HTTPServer struct {
	Address     string        `default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}
type ClientsConfig struct {
	Accounts Client `yaml:"accounts"`
	Lessons  Client `yaml:"lessons"`
}

func main() {

	// Загружаем конфиги
	cfg := Config{}

	config.MustConfig(&cfg)

	// Настраиваем логер
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	).With(
		slog.String("address", cfg.Address),
	)

	// Информируем о том что настроили
	logger.Info("server initialized")

	// Инициализирем grpc клиент c аккаунтами
	accountService, err := accounts.NewAccountClient(
		context.Background(),
		logger,
		cfg.Clients.Accounts.Address,
		cfg.Clients.Accounts.Timeout,
		cfg.Clients.Accounts.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to init account client", lib.Err(err))
		os.Exit(1)
	}
	// Инициализирем grpc клиент c аккаунтами
	lessonService, err := lessons.NewLessonClient(
		context.Background(),
		logger,
		cfg.Clients.Lessons.Address,
		cfg.Clients.Lessons.Timeout,
		cfg.Clients.Lessons.RetriesCount,
	)
	if err != nil {
		logger.Error("failed to init lessons client", lib.Err(err))
		os.Exit(1)
	}

	// Инициализируем сервисный слой и слой хендлера
	serviceLayer := service.NewService(accountService, lessonService)
	handlerLayer := handler.NewHandler(serviceLayer, logger)
	srv := server.Server{}

	go func() {
		if err := srv.Run(cfg.Address, handlerLayer.InitRoutes()); err != nil {
			logger.Error("error occured while running http server %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	if err := srv.ShutDown(context.Background()); err != nil {
		logger.Error("Error while stopping")
	}
}
