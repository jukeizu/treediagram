package scheduling

import (
	pb "github.com/jukeizu/treediagram/api/scheduling"
	mdb "github.com/shawntoffel/GoMongoDb"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName   = "scheduling"
	CollectionName = "jobs"
)

type JobStorage interface {
	mdb.Storage
	Create(*pb.Job) error
	Jobs(*pb.Schedule) ([]*pb.Job, error)
	Disable(id string) error
}

type jobStorage struct {
	mdb.Store
}

func NewJobStorage(url string) (JobStorage, error) {
	c := mdb.DbConfig{
		Url:            url,
		DatabaseName:   DatabaseName,
		CollectionName: CollectionName,
	}

	store, err := mdb.NewStorage(c)

	j := jobStorage{}
	j.Session = store.Session
	j.Collection = store.Collection

	return &j, err
}

func (store *jobStorage) Create(job *pb.Job) error {
	return store.Collection.Insert(job)
}

func (store *jobStorage) Jobs(schedule *pb.Schedule) ([]*pb.Job, error) {
	jobs := []*pb.Job{}

	if schedule == nil {
		err := store.Collection.Find(bson.M{"enabled": true}).All(&jobs)
		return jobs, err
	}

	scheduleSelector := bson.M{
		"$and": []bson.M{
			buildTimeSelector("schedule.minute", schedule.Minute),
			buildTimeSelector("schedule.hour", schedule.Hour),
			buildTimeSelector("schedule.dayofmonth", schedule.DayOfMonth),
			buildTimeSelector("schedule.month", schedule.Month),
			buildTimeSelector("schedule.dayofweek", schedule.DayOfWeek),
			bson.M{"enabled": true},
		},
	}

	err := store.Collection.Find(scheduleSelector).All(&jobs)

	return jobs, err
}

func (store *jobStorage) Disable(id string) error {
	_, err := store.Collection.Upsert(bson.M{"id": id}, bson.M{"$set": bson.M{"enabled": false}})

	return err
}

func buildTimeSelector(field string, value string) bson.M {
	return bson.M{field: bson.M{"$in": []string{value, ""}}}
}
