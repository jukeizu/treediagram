package startup

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shawntoffel/gossage"
)

type MigrationRunner struct {
	Logger  zerolog.Logger
	Storage *Storage
}

func NewMigrationRunner(logger zerolog.Logger, dbUrl string) (*MigrationRunner, error) {
	logger = logger.With().Str("component", "migrator").Logger()

	gossage.Logger = func(format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a...)
		logger.Info().Msg(msg)
	}

	storage, err := NewStorage(logger, dbUrl)
	if err != nil {
		return nil, errors.New("db: " + err.Error())
	}

	return &MigrationRunner{Storage: storage, Logger: logger}, nil
}

func (m *MigrationRunner) Migrate() error {
	m.Logger.Info().Msg("starting migrations")

	err := m.Storage.Migrate()
	if err != nil {
		return err
	}

	m.Logger.Info().Msg("migrations complete")

	return nil
}
