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
	IntentStorage    intent.IntentStorage
	JobStorage       scheduler.JobStorage
	UserStorage      user.UserStorage
}

func NewStorage(logger zerolog.Logger, mdbUrl string, dbUrl string) (*Storage, error) {
	gooseLogger := goosezerolog.New(logger)
	goose.SetLogger(gooseLogger)
	goose.SetDialect("postgres")

	processorStorage, err := processor.NewStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	commandStorage, err := intent.NewIntentStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	jobStorage, err := scheduler.NewJobStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	userStorage, err := user.NewUserStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	err = userStorage.Migrate()
	if err != nil {
		return nil, err
	}

	s := &Storage{
		ProcessorStorage: processorStorage,
		IntentStorage:    commandStorage,
		JobStorage:       jobStorage,
		UserStorage:      userStorage,
	}

	return s, nil
}

func (s *Storage) Close() {
	s.ProcessorStorage.Close()
	s.IntentStorage.Close()
	s.JobStorage.Close()
}
