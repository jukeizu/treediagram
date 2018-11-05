package command

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "processor"
	CollectionName = "reply"
)

type Reply struct {
	Id              bson.ObjectId `bson:"_id,omitempty"`
	ProcessingReply processing.Reply
	Errors          []string
}

type Storage interface {
	mdb.Storage
	SaveReply(Reply) (string, error)
	Reply(id string) (Reply, error)
}

type storage struct {
	mdb.Store
}

func NewStorage(url string) (Storage, error) {
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

func (s *storage) SaveReply(r Reply) (string, error) {
	r.Id = bson.NewObjectId()

	err := s.Collection.Insert(r)

	return r.Id.Hex(), err
}

func (s *storage) Reply(id string) (Reply, error) {
	r := Reply{}

	if !bson.IsObjectIdHex(id) {
		return r, nil
	}

	err := s.Collection.FindId(bson.ObjectIdHex(id)).One(&r)

	return r, err
}
