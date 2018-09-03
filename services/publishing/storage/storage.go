package storage

import (
	pb "github.com/jukeizu/treediagram/api/publishing"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

type MessageStorage interface {
	mdb.Storage
	Save(*pb.Message) error
	Message(id string) (*pb.Message, error)
}

type messageStorage struct {
	mdb.Store
}

func NewMessageStorage(dbConfig mdb.DbConfig) (MessageStorage, error) {
	store, err := mdb.NewStorage(dbConfig)

	m := messageStorage{}
	m.Session = store.Session
	m.Collection = store.Collection

	return &m, err
}

func (store *messageStorage) Save(message *pb.Message) error {
	return store.Collection.Insert(message)
}

func (store *messageStorage) Message(id string) (*pb.Message, error) {
	message := pb.Message{}

	err := store.Collection.Find(bson.M{"id": id}).One(&message)

	return &message, err
}
