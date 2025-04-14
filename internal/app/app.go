package app

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/magneless/pvz/internal/repository/postgresql"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/magneless/pvz/pkg/config"
	database "github.com/magneless/pvz/pkg/database/postgres"
	"github.com/magneless/pvz/pkg/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Application struct {
	Config            *config.Config
	DB                *sql.DB
	PVZUsecase        *usecase.PVZUsecase
	ReceptionUsecase  *usecase.ReceptionUsecase
	DummyLoginUsecase *usecase.DummyLoginUsecase
	ProductUsecase    *usecase.ProductUsecase
	Logger            *slog.Logger
}

func New() *Application {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)

	db, err := database.NewPostgres(cfg.Storage)
	if err != nil {
		logger.Error("failed to init or connnect to db", sl.Err(err))
		os.Exit(1)
	}

	pvzRepo := postgresql.NewPVZRepository(db)
	pvzUsecase := usecase.NewPVZUsecase(pvzRepo)

	receptionRepo := postgresql.NewReceptionRepository(db)
	receptionUsecase := usecase.NewReceptionUsecase(receptionRepo)

	dummyLoginUsecase := usecase.NewDumyyLoginUsecase()

	productRepository := postgresql.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepository)

	return &Application{
		Config:            cfg,
		DB:                db,
		PVZUsecase:        pvzUsecase,
		ReceptionUsecase:  receptionUsecase,
		DummyLoginUsecase: dummyLoginUsecase,
		ProductUsecase:    productUsecase,
		Logger:            logger,
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		// fallback
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
