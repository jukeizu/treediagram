package main

import (
	"github.com/jukeizu/treediagram/publishing"
	"github.com/jukeizu/treediagram/registration"
	"github.com/jukeizu/treediagram/scheduling"
	"github.com/jukeizu/treediagram/user"
)

type Storage struct {
	MessageStorage publishing.MessageStorage
	CommandStorage registration.CommandStorage
	JobStorage     scheduling.JobStorage
	UserStorage    user.UserStorage
}

func NewStorage(dbUrl string) (*Storage, error) {

	messageStorage, err := publishing.NewMessageStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	commandStorage, err := registration.NewCommandStorage(dbUrl)
	if err != nil {
		return nil, err
	}

	jobStorage, err := scheduling.NewJobStorage(dbUrl)
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
