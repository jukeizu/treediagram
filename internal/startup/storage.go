package startup

import (
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	"github.com/rs/zerolog"
)

type Storage struct {
	ProcessorStorage    processor.Storage
	IntentRepository    intent.Repository
	SchedulerRepository scheduler.Repository
	UserRepository      user.Repository
}

func NewStorage(logger zerolog.Logger, mdbUrl string, dbUrl string) (*Storage, error) {
	processorStorage, err := processor.NewStorage(mdbUrl)
	if err != nil {
		return nil, err
	}

	intentRepository, err := intent.NewRepository(dbUrl)
	if err != nil {
		return nil, err
	}

	schedulerRepository, err := scheduler.NewRepository(dbUrl)
	if err != nil {
		return nil, err
	}

	userRepository, err := user.NewRepository(dbUrl)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		ProcessorStorage:    processorStorage,
		IntentRepository:    intentRepository,
		SchedulerRepository: schedulerRepository,
		UserRepository:      userRepository,
	}

	return s, nil
}

func (s *Storage) Migrate() error {
	err := s.UserRepository.Migrate()
	if err != nil {
		return err
	}

	err = s.IntentRepository.Migrate()
	if err != nil {
		return err
	}

	err = s.SchedulerRepository.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() {
	s.ProcessorStorage.Close()
}
