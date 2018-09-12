package user

import (
	pb "github.com/jukeizu/treediagram/api/user"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "user"
	CollectionName = "preferences"
)

type UserStorage interface {
	mdb.Storage

	Preference(*pb.PreferenceRequest) (*pb.Preference, error)
	SetServer(*pb.SetServerRequest) (*pb.Preference, error)
}

type storage struct {
	mdb.Store
}

func NewUserStorage(url string) (UserStorage, error) {
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

func (s *storage) Preference(req *pb.PreferenceRequest) (*pb.Preference, error) {
	preference := &pb.Preference{}

	err := s.Collection.Find(bson.M{"userid": req.UserId}).One(&preference)

	return preference, err
}

func (s *storage) SetServer(req *pb.SetServerRequest) (*pb.Preference, error) {
	_, err := s.Collection.Upsert(bson.M{"userid": req.UserId}, bson.M{"$set": bson.M{"serverid": req.ServerId}})

	if err != nil {
		return nil, err
	}

	return s.Preference(&pb.PreferenceRequest{UserId: req.UserId})
}
