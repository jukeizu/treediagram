package registration

type Service interface {
	Add(command Command) (Command, error)
	Disable(id string) error
	Query(query CommandQuery) (CommandQueryResult, error)
}

type service struct {
	CommandStorage CommandStorage
}

func NewService(commandStorage CommandStorage) (Service, error) {
	return &service{CommandStorage: commandStorage}, nil
}

func (s *service) Add(command Command) (Command, error) {
	err := s.CommandStorage.Save(command)

	return command, err
}

func (s *service) Disable(id string) error {
	return s.CommandStorage.Disable(id)
}

func (s *service) Query(query CommandQuery) (CommandQueryResult, error) {
	if query.PageSize < 1 {
		query.PageSize = 50
	}

	result := CommandQueryResult{}

	commands, err := s.CommandStorage.Query(query)

	if err != nil {
		return result, err
	}

	result.Commands = commands

	if len(result.Commands) == query.PageSize {
		result.HasMore = true
	}

	return result, err
}
