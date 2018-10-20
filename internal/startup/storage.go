package startup

import (
	"github.com/jukeizu/treediagram/processor/user"
	"github.com/jukeizu/treediagram/publisher"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
)

type Storage struct {
	MessageStorage publisher.MessageStorage
	CommandStorage registry.CommandStorage
	JobStorage     scheduler.JobStorage
	UserStorage    user.UserStorage
}

func NewStorage(dbUrl string) (*Storage, error) {

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
		MessageStorage: messageStorage,
		CommandStorage: commandStorage,
		JobStorage:     jobStorage,
		UserStorage:    userStorage,
	}

	return s, nil
}

func (s *Storage) Close() {
	s.MessageStorage.Close()
	s.CommandStorage.Close()
	s.JobStorage.Close()
	s.UserStorage.Close()
}
