package registry

import (
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/registration"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "registration"
	CollectionName = "commands"
)

type Command struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Server   string        `json:"server"`
	Name     string        `json:"name"`
	Regex    string        `json:"regex"`
	Endpoint string        `json:"endpoint"`
	Help     string        `json:"help"`
	Enabled  bool          `json:"enabled"`
}

type CommandStorage interface {
	mdb.Storage

	Save(pb.Command) error
	Disable(string) error
	Query(pb.QueryCommandsRequest) ([]*pb.Command, error)
}

type storage struct {
	mdb.Store
}

func NewCommandStorage(url string) (CommandStorage, error) {
	c := mdb.DbConfig{
		Url:            url,
		DatabaseName:   DatabaseName,
		CollectionName: CollectionName,
	}

	store, err := mdb.NewStorage(c)
	if err != nil {
		return nil, err
	}

	j := storage{}
	j.Session = store.Session
	j.Collection = store.Collection

	return &j, err
}

func (s *storage) Save(pbCommand pb.Command) error {
	c := Command{
		Server:   pbCommand.Server,
		Name:     pbCommand.Name,
		Regex:    pbCommand.Regex,
		Endpoint: pbCommand.Endpoint,
		Help:     pbCommand.Help,
		Enabled:  pbCommand.Enabled,
	}

	return s.Collection.Insert(c)
}

func (s *storage) Disable(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("db error: the following id is invalid: %s", id)
	}
	_, err := s.Collection.Upsert(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"enabled": false}})
	if err != nil {
		return fmt.Errorf("db error: %s", err)
	}

	return nil
}

func (s *storage) Query(query pb.QueryCommandsRequest) ([]*pb.Command, error) {
	commands := []Command{}

	bsonQuery := []bson.M{
		bson.M{"server": bson.M{"$in": []string{query.Server, ""}}},
		bson.M{"enabled": true},
	}
	if bson.IsObjectIdHex(query.LastId) {
		bsonQuery = append(bsonQuery, bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex(query.LastId)}})
	}

	pbCommands := []*pb.Command{}

	err := s.Collection.Find(bson.M{"$and": bsonQuery}).Limit(int(query.PageSize)).All(&commands)
	if err != nil {
		return pbCommands, fmt.Errorf("db error: %s", err)
	}

	for _, command := range commands {
		mapped := &pb.Command{
			Id:       command.Id.Hex(),
			Server:   command.Server,
			Name:     command.Name,
			Regex:    command.Regex,
			Endpoint: command.Endpoint,
			Help:     command.Help,
			Enabled:  command.Enabled,
		}
		pbCommands = append(pbCommands, mapped)
	}

	return pbCommands, nil
}
