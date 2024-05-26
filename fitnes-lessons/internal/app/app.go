package app

import (
	"context"
	"fitnes-lessons/internal/app/grpcapp"
	"fitnes-lessons/internal/repository/postgres"
	"fitnes-lessons/internal/service"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

type Config struct {
	Repo postgres.Config
	GRPC grpcapp.Config
}

func New(ctx context.Context, log *slog.Logger, cfg *Config) *App {
	storage, err := postgres.NewPostgresRepository(&cfg.Repo)
	if err != nil {
		panic(err)
	}

	ls, err := service.NewLessonService(ctx, storage, storage, storage)
	if err != nil {
		panic(err)
	}

	grpcApp := grpcapp.NewGrpcApp(log, ls, &cfg.GRPC)

	return &App{
		GRPCServer: grpcApp,
	}
}
