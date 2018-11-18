package startup

import (
	"github.com/jukeizu/treediagram/publisher"
	"github.com/jukeizu/treediagram/registry"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
)

type Storage struct {
	MessageStorage publisher.MessageStorage
	IntentStorage  registry.IntentStorage
	JobStorage     scheduler.JobStorage
	UserStorage    user.UserStorage
}

func NewStorage(dbUrl string) (*Storage, error) {
	messageStorage, err := publisher.NewMessageStorage(dbUrl)
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
		MessageStorage: messageStorage,
		IntentStorage:  commandStorage,
		JobStorage:     jobStorage,
		UserStorage:    userStorage,
	}

	return s, nil
}

func (s *Storage) Close() {
	s.MessageStorage.Close()
	s.IntentStorage.Close()
	s.JobStorage.Close()
	s.UserStorage.Close()
}
