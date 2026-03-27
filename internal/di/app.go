package di

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/SilverName608/Effective-Mobile/internal/api"
	"github.com/SilverName608/Effective-Mobile/internal/config"
	"github.com/SilverName608/Effective-Mobile/internal/database"
	"github.com/SilverName608/Effective-Mobile/internal/repository"
	"github.com/SilverName608/Effective-Mobile/internal/service"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(database.NewPostgres),
		fx.Provide(newLogger),

		fx.Provide(fx.Annotate(
			repository.NewSubscriptionRepository,
			fx.As(new(repository.SubscriptionRepositoryI)),
		)),

		fx.Provide(fx.Annotate(
			service.NewSubscriptionService,
			fx.As(new(service.SubscriptionServiceI)),
		)),

		fx.Provide(api.NewSubscriptionHandler),
		fx.Provide(api.NewRouter),

		fx.Invoke(runServer),
	)
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func runServer(cfg *config.Config, router chi.Router, db *gorm.DB, log *logrus.Logger) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := runMigrations(sqlDB); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Infof("server starting -> http://localhost%s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	log.Println("migrations applied")
	return nil
}
