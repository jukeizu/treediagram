package storage

import (
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

type MessageStorage interface {
	mdb.Storage
	SaveMessage(Message) error
	GetMessage(id string) (Message, error)
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

func (store *messageStorage) SaveMessage(message Message) error {
	return store.Collection.Insert(message)
}

func (store *messageStorage) GetMessage(id string) (Message, error) {
	message := Message{}

	err := store.Collection.Find(bson.M{"id": id}).One(&message)

	return message, err
}
