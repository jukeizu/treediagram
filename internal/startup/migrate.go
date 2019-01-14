package startup

import (
	"errors"

	"github.com/rs/zerolog"
)

type MigrationRunner struct {
	Logger  zerolog.Logger
	Storage *Storage
}

func NewMigrationRunner(logger zerolog.Logger, mdbUrl string, dbUrl string) (*MigrationRunner, error) {
	logger = logger.With().Str("component", "migrator").Logger()

	storage, err := NewStorage(logger, mdbUrl, dbUrl)
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
