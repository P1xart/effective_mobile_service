package migrations

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/P1xart/effective_mobile_service/internal/config"
	"github.com/P1xart/effective_mobile_service/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(log *slog.Logger, fs *embed.FS, cfg *config.Postgresql) {
	log = log.With(slog.String("component", "migrations"))

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	source, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Error("failed to read migrations source", logger.Error(err))
		return
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, makeMigrateUrl(dbUrl))
	if err != nil {
		log.Error("failed to initialization the migrations instance", logger.Error(err))
		return
	}

	err = instance.Up()

	switch err {
	case nil:
		log.Debug("the migration schema successfully upgraded!")
	case migrate.ErrNoChange:
		log.Debug("the migration schema not changed")
	default:
		log.Error("could not apply the migration schema", logger.Error(err))
		os.Exit(1)
	}
}
