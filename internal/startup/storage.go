package startup

import (
	"github.com/jukeizu/treediagram/processor/command"
	"github.com/jukeizu/treediagram/publisher"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
)

type Storage struct {
	ProcessorStorage command.Storage
	MessageStorage   publisher.MessageStorage
	CommandStorage   registry.CommandStorage
	JobStorage       scheduler.JobStorage
	UserStorage      user.UserStorage
}

func NewStorage(dbUrl string) (*Storage, error) {
	processorStorage, err := command.NewStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	messageStorage, err := publisher.NewMessageStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	commandStorage, err := registry.NewCommandStorage(dbUrl)
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
		MessageStorage:   messageStorage,
		CommandStorage:   commandStorage,
		JobStorage:       jobStorage,
		UserStorage:      userStorage,
	}

	return s, nil
}

func (s *Storage) Close() {
	s.ProcessorStorage.Close()
	s.MessageStorage.Close()
	s.CommandStorage.Close()
	s.JobStorage.Close()
	s.UserStorage.Close()
}
