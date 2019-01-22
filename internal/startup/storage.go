package startup

import (
	"github.com/jukeizu/treediagram/intent"
	"github.com/jukeizu/treediagram/processor"
	"github.com/jukeizu/treediagram/scheduler"
	"github.com/jukeizu/treediagram/user"
	"github.com/rs/zerolog"
)

type Storage struct {
	ProcessorRepository processor.Repository
	IntentRepository    intent.Repository
	SchedulerRepository scheduler.Repository
	UserRepository      user.Repository
}

func NewStorage(logger zerolog.Logger, dbUrl string) (*Storage, error) {
	processorRepository, err := processor.NewRepository(dbUrl)
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
		ProcessorRepository: processorRepository,
		IntentRepository:    intentRepository,
		SchedulerRepository: schedulerRepository,
		UserRepository:      userRepository,
	}

	return s, nil
}

func (s *Storage) Migrate() error {
	err := s.ProcessorRepository.Migrate()
	if err != nil {
		return err
	}

	err = s.UserRepository.Migrate()
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
