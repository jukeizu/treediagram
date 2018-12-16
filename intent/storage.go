package intent

import (
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/intent"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "registration"
	CollectionName = "intents"
)

type Intent struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Server   string        `json:"server"`
	Name     string        `json:"name"`
	Regex    string        `json:"regex"`
	Mention  bool          `json:"mention"`
	Response string        `json:"response"`
	Endpoint string        `json:"endpoint"`
	Help     string        `json:"help"`
	Enabled  bool          `json:"enabled"`
}

type IntentStorage interface {
	mdb.Storage

	Save(pb.Intent) error
	Disable(string) error
	Query(pb.QueryIntentsRequest) ([]*pb.Intent, error)
}

type storage struct {
	mdb.Store
}

func NewIntentStorage(url string) (IntentStorage, error) {
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

func (s *storage) Save(pbIntent pb.Intent) error {
	i := Intent{
		Server:   pbIntent.Server,
		Name:     pbIntent.Name,
		Regex:    pbIntent.Regex,
		Mention:  pbIntent.Mention,
		Response: pbIntent.Response,
		Endpoint: pbIntent.Endpoint,
		Help:     pbIntent.Help,
		Enabled:  pbIntent.Enabled,
	}

	return s.Collection.Insert(i)
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

func (s *storage) Query(query pb.QueryIntentsRequest) ([]*pb.Intent, error) {
	intents := []Intent{}

	bsonQuery := []bson.M{
		bson.M{"server": bson.M{"$in": []string{query.Server, ""}}},
		bson.M{"enabled": true},
	}

	pbIntents := []*pb.Intent{}

	err := s.Collection.Find(bson.M{"$and": bsonQuery}).All(&intents)
	if err != nil {
		return pbIntents, fmt.Errorf("db error: %s", err)
	}

	for _, i := range intents {
		mapped := &pb.Intent{
			Id:       i.Id.Hex(),
			Server:   i.Server,
			Name:     i.Name,
			Regex:    i.Regex,
			Mention:  i.Mention,
			Response: i.Response,
			Endpoint: i.Endpoint,
			Help:     i.Help,
			Enabled:  i.Enabled,
		}
		pbIntents = append(pbIntents, mapped)
	}

	return pbIntents, nil
}
