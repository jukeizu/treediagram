package storage

import (
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

type MessageRequestStorage interface {
	mdb.Storage
	SaveMessageRequest(MessageRequest) error
	GetMessageRequest(id string) (MessageRequest, error)
}

type messageRequestStorage struct {
	mdb.Store
}

func NewMessageRequestStorage(dbConfig mdb.DbConfig) (MessageRequestStorage, error) {
	store, err := mdb.NewStorage(dbConfig)

	m := messageRequestStorage{}
	m.Session = store.Session
	m.Collection = store.Collection

	return &m, err
}

func (store *messageRequestStorage) SaveMessageRequest(messageRequest MessageRequest) error {
	return store.Collection.Insert(messageRequest)
}

func (store *messageRequestStorage) GetMessageRequest(id string) (MessageRequest, error) {
	messageRequest := MessageRequest{}

	err := store.Collection.Find(bson.M{"id": id}).One(&messageRequest)

	return messageRequest, err
}
