package command

import (
	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "processor"
	CollectionName = "commands"
)

type Storage interface {
	mdb.Storage
	SaveMatch(Match) (string, error)
	Match(id string) (Match, error)
}

type storage struct {
	mdb.Store
}

type Match struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Request Request
	Command Command
}

type Request struct {
	Id            string
	Source        string
	CorrelationId string
	Bot           User
	Author        User
	ChannelId     string
	ServerId      string
	Mentions      []User
	Content       string
}

type User struct {
	Id   string
	Name string
}

type Command struct {
	Id       string
	Endpoint string
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

func (s *storage) SaveMatch(m Match) (string, error) {
	m.Id = bson.NewObjectId()

	err := s.Collection.Insert(m)

	return m.Id.Hex(), err
}

func (s *storage) Match(id string) (Match, error) {
	m := Match{}

	if !bson.IsObjectIdHex(id) {
		return m, nil
	}

	err := s.Collection.FindId(bson.ObjectIdHex(id)).One(&m)

	return m, err
}

func toRequest(pr processing.TreediagramRequest) Request {
	r := Request{
		Source:        pr.Source,
		CorrelationId: pr.CorrelationId,
		Bot:           toUser(pr.Bot),
		Author:        toUser(pr.Author),
		ChannelId:     pr.ChannelId,
		ServerId:      pr.ServerId,
		Content:       pr.Content,
	}

	for _, m := range pr.Mentions {
		r.Mentions = append(r.Mentions, toUser(m))
	}

	return r
}

func toUser(pu *processing.User) User {
	u := User{
		Id:   pu.Id,
		Name: pu.Name,
	}

	return u
}
