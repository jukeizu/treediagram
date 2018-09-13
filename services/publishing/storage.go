package publishing

import (
	pb "github.com/jukeizu/treediagram/api/publishing"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "publishing"
	CollectionName = "messages"
)

type MessageStorage interface {
	mdb.Storage
	Save(*pb.Message) error
	Message(id string) (*pb.Message, error)
}

type messageStorage struct {
	mdb.Store
}

func NewMessageStorage(url string) (MessageStorage, error) {
	c := mdb.DbConfig{
		Url:            url,
		DatabaseName:   DatabaseName,
		CollectionName: CollectionName,
	}

	store, err := mdb.NewStorage(c)
	if err != nil {
		return nil, err
	}

	m := messageStorage{}
	m.Session = store.Session
	m.Collection = store.Collection

	return &m, err
}

func (store *messageStorage) Save(message *pb.Message) error {
	return store.Collection.Insert(message)
}

func (store *messageStorage) Message(id string) (*pb.Message, error) {
	message := &pb.Message{}

	err := store.Collection.Find(bson.M{"id": id}).One(message)

	return message, err
}
