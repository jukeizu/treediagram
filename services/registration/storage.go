package registration

import (
	"errors"

	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "registration"
	CollectionName = "commands"
)

type CommandStorage interface {
	mdb.Storage

	Save(Command) error
	Disable(string) error
	Query(CommandQuery) ([]Command, error)
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

	j := storage{}
	j.Session = store.Session
	j.Collection = store.Collection

	return &j, err
}

func (s *storage) Save(c Command) error {
	return s.Collection.Insert(c)
}

func (s *storage) Disable(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("The following id is invalid: " + id)
	}
	_, err := s.Collection.Upsert(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"enabled": false}})

	return err
}

func (s *storage) Query(query CommandQuery) ([]Command, error) {
	commands := []Command{}

	bsonQuery := []bson.M{
		bson.M{"server": query.Server},
		bson.M{"enabled": true},
	}
	if bson.IsObjectIdHex(query.LastId) {
		bsonQuery = append(bsonQuery, bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex(query.LastId)}})
	}

	err := s.Collection.Find(bson.M{"$and": bsonQuery}).Limit(query.PageSize).All(&commands)

	return commands, err
}
