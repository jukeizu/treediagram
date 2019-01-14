package startup

import (
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/pkg/goosezerolog"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	"github.com/pressly/goose"
	"github.com/rs/zerolog"
)

type Storage struct {
	ProcessorStorage processor.Storage
	IntentDb         intent.IntentDb
	JobStorage       scheduler.JobStorage
	UserDb           user.UserDb
}

func NewStorage(logger zerolog.Logger, mdbUrl string, dbUrl string) (*Storage, error) {
	gooseLogger := goosezerolog.New(logger)
	goose.SetLogger(gooseLogger)
	goose.SetDialect("postgres")

	processorStorage, err := processor.NewStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	intentDb, err := intent.NewIntentDb(dbUrl)
	if err != nil {
		return nil, err
	}

	jobStorage, err := scheduler.NewJobStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	userDb, err := user.NewUserDb(dbUrl)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		ProcessorStorage: processorStorage,
		IntentDb:         intentDb,
		JobStorage:       jobStorage,
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

	return nil
}

func (s *Storage) Close() {
	s.ProcessorStorage.Close()
	s.JobStorage.Close()
}
