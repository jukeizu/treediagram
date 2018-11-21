package startup

import (
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
)

type Storage struct {
	ProcessorStorage processor.Storage
	IntentStorage    registry.IntentStorage
	JobStorage       scheduler.JobStorage
	UserStorage      user.UserStorage
}

func NewStorage(dbUrl string) (*Storage, error) {
	processorStorage, err := processor.NewStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	commandStorage, err := registry.NewIntentStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	jobStorage, err := scheduler.NewJobStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	userStorage, err := user.NewUserStorage(dbUrl)
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
	s.UserStorage.Close()
}
