package registration

type Service interface {
	Add(command Command) (Command, error)
	Remove(id string) error
	Commands() ([]Command, error)
}

type service struct {
}

func NewService() (Service, error) {
	return &service{}, nil
}

func (s *service) Add(command Command) (Command, error) {
	return Command{}, nil
}

func (s *service) Remove(id string) error {
	return nil
}

func (s *service) Commands() ([]Command, error) {
	return []Command{}, nil
}
