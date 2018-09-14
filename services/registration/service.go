package registration

import "gopkg.in/mgo.v2/bson"

type Command struct {
	Id             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Server         string        `json:"server"`
	Name           string        `json:"name"`
	Regex          string        `json:"regex"`
	RequireMention bool          `json:"requireMention"`
	Endpoint       string        `json:"endpoint"`
	Help           string        `json:"help"`
	Enabled        bool          `json:"enabled"`
}

type CommandQuery struct {
	Server   string `json:"server"`
	LastId   string `json:"lastId"`
	PageSize int    `json:"pageSize"`
}

type CommandQueryResult struct {
	Commands []Command `json:"commands"`
	HasMore  bool      `json:"hasMore"`
}

type Service interface {
	Add(command Command) (Command, error)
	Disable(id string) error
	Query(query CommandQuery) (CommandQueryResult, error)
}

type service struct {
	CommandStorage CommandStorage
}

func NewService(commandStorage CommandStorage) Service {
	return &service{CommandStorage: commandStorage}
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
