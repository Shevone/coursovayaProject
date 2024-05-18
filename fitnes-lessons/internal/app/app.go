package app

import (
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

func New(log *slog.Logger, cfg *Config) *App {
	storage, err := postgres.NewPostgresRepository(&cfg.Repo)
	if err != nil {
		panic(err)
	}

	lessonService := service.NewLessonService(log, storage, storage, storage)

	grpcApp := grpcapp.NewGrpcApp(log, lessonService, &cfg.GRPC)

	return &App{
		GRPCServer: grpcApp,
	}
}
