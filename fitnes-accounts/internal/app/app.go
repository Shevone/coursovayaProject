package app

import (
	"fitnes-account/internal/app/grpcapp"
	"fitnes-account/internal/repository/postgres"
	"fitnes-account/internal/service"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

type Config struct {
	Repo    postgres.Config
	GRPC    grpcapp.Config
	Service service.Config
}

func New(
	log *slog.Logger,
	cfg *Config,
) *App {
	storage, err := postgres.NewPostgresRepository(&cfg.Repo)
	if err != nil {
		panic(err)
	}

	authService := service.NewAccountService(log, storage, storage, &cfg.Service)

	grpcApp := grpcapp.NewGrpcApp(log, authService, &cfg.GRPC)

	return &App{
		GRPCServer: grpcApp,
	}
}
