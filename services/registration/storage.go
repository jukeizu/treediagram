package registration

import (
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
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

func NewCommandStorage(dbConfig mdb.DbConfig) (CommandStorage, error) {
	store, err := mdb.NewStorage(dbConfig)

	j := storage{}
	j.Session = store.Session
	j.Collection = store.Collection

	return &j, err
}

func (s *storage) Save(c Command) error {
	return s.Collection.Insert(c)
}

func (s *storage) Disable(id string) error {
	_, err := s.Collection.Upsert(bson.M{"_id": id}, bson.M{"$set": bson.M{"enabled": false}})

	return err
}

func (s *storage) Query(query CommandQuery) ([]Command, error) {
	commands := []Command{}

	err := s.Collection.Find(bson.M{"_id": bson.M{"$gt": query.LastId}, "server": query.Server}).Limit(query.PageSize).All(&commands)

	return commands, err
}
