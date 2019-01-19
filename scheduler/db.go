package scheduler

import (
	"database/sql"
	"fmt"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/scheduling"
	migration "github.com/jukeizu/treediagram/scheduler/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_scheduler"
)

type JobDb interface {
	Create(*pb.Job) error
	Jobs(*pb.Schedule) ([]*pb.Job, error)
	Disable(id string) error
	Migrate() error
}

type jobDb struct {
	Db *sql.DB
}

func NewJobDb(url string) (JobDb, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	j := jobDb{
		Db: db,
	}

	return &j, nil
}

func (j *jobDb) Migrate() error {
	g, err := gossage.New(j.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(migration.CreateTableJob20190119052738{})
	if err != nil {
		return err
	}

	return g.Up()
}

func (j *jobDb) Create(job *pb.Job) error {
	return nil
}

func (j *jobDb) Jobs(schedule *pb.Schedule) ([]*pb.Job, error) {
	jobs := []*pb.Job{}

	return jobs, nil
}

func (j *jobDb) Disable(id string) error {
	return nil
}
