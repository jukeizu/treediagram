package processor

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName               = "processor"
	CommandCollectionName      = "commands"
	CommandEventCollectionName = "commandevents"
	ReplyCollectionName        = "replies"
)

type Storage struct {
	Session                *mgo.Session
	CommandCollection      *mgo.Collection
	CommandEventCollection *mgo.Collection
	ReplyCollection        *mgo.Collection
}

func NewStorage(url string) (Storage, error) {
	s := Storage{}
	session, err := mgo.Dial(url)
	if err != nil {
		return s, err
	}

	s.Session = session

	db := s.Session.DB(DatabaseName)

	s.CommandCollection = db.C(CommandCollectionName)
	s.CommandEventCollection = db.C(CommandEventCollectionName)
	s.ReplyCollection = db.C(ReplyCollectionName)

	return s, nil
}

func (s Storage) Close() {
	if s.Session == nil {
		return
	}

	s.Session.Close()
}

func (s Storage) SaveCommand(c Command) error {
	return s.CommandCollection.Insert(c)
}

func (s Storage) Command(id string) (Command, error) {
	c := Command{}

	err := s.CommandCollection.Find(bson.M{"id": id}).One(&c)

	return c, err
}

func (s Storage) SaveCommandEvent(e CommandEvent) error {
	return s.CommandEventCollection.Insert(e)
}

func (s Storage) CommandEvents(id string) ([]CommandEvent, error) {
	c := []CommandEvent{}

	err := s.CommandEventCollection.Find(bson.M{"commandid": id}).All(&c)

	return c, err
}

func (s Storage) SaveReply(r processing.Reply) error {
	return s.ReplyCollection.Insert(r)
}

func (s Storage) Reply(id string) (*processing.Reply, error) {
	r := &processing.Reply{}

	err := s.ReplyCollection.Find(bson.M{"id": id}).One(r)

	return r, err
}
