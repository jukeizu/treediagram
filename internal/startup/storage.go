package startup

import (
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	"github.com/rs/zerolog"
)

type Storage struct {
	ProcessorStorage processor.Storage
	IntentDb         intent.IntentDb
	JobDb            scheduler.JobDb
	UserDb           user.UserDb
}

func NewStorage(logger zerolog.Logger, mdbUrl string, dbUrl string) (*Storage, error) {
	processorStorage, err := processor.NewStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	intentDb, err := intent.NewIntentDb(dbUrl)
	if err != nil {
		return nil, err
	}

	jobDb, err := scheduler.NewJobDb(dbUrl)
	if err != nil {
		return nil, err
	}

	userDb, err := user.NewUserDb(logger, dbUrl)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		ProcessorStorage: processorStorage,
		IntentDb:         intentDb,
		JobDb:            jobDb,
		UserDb:           userDb,
	}

	return s, nil
}

func (s *Storage) Migrate() error {
	err := s.UserDb.Migrate()
	if err != nil {
		return err
	}

	err = s.IntentDb.Migrate()
	if err != nil {
		return err
	}

	err = s.JobDb.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() {
	s.ProcessorStorage.Close()
}
