package app

import (
	"fitnes-account/internal/app/grpcapp"
	"fitnes-account/internal/repository/sqlite"
	"fitnes-account/internal/service"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.NewSqliteRepository(storagePath)
	if err != nil {
		panic(err)
	}

	authService := service.NewAccountService(log, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewGrpcApp(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
