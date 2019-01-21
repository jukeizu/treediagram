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
	Create(*pb.Job) (*pb.Job, error)
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

func (j *jobDb) Create(job *pb.Job) (*pb.Job, error) {
	q := `INSERT INTO job (
			type,
			content,
			userId,
			destination,
			minute,
			hour,
			dayOfMonth,
			month,
			dayOfWeek,	
			year,
			enabled
		) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id,
			type,
			content,
			userId,
			destination,
			minute,
			hour,
			dayOfMonth,
			month,
			dayOfWeek,	
			year,
			enabled,
			created::INT`

	job, err := j.queryJob(q,
		job.Type,
		job.Content,
		job.UserId,
		job.Destination,
		job.Schedule.Minute,
		job.Schedule.Hour,
		job.Schedule.DayOfMonth,
		job.Schedule.Month,
		job.Schedule.DayOfWeek,
		job.Schedule.Year,
		job.Enabled,
	)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j *jobDb) Jobs(schedule *pb.Schedule) ([]*pb.Job, error) {
	if schedule == nil {
		return j.allJobs()
	}

	q := `SELECT id,
		type,
		content,
		userId,
		destination,
		minute,
		hour,
		dayOfMonth,
		month,
		dayOfWeek,	
		year,
		enabled,
		created::INT
	FROM job
	WHERE minute IN($1, '')
	AND hour IN($2, '')
	AND dayOfMonth IN($3, '')
	AND month IN($4, '')
	AND dayOfWeek IN ($5, '')
	AND year IN ($6, '')
	AND enabled = true`

	jobs, err := j.queryJobs(q,
		schedule.Minute,
		schedule.Hour,
		schedule.DayOfMonth,
		schedule.Month,
		schedule.DayOfWeek,
		schedule.Year,
	)

	return jobs, err
}

func (j *jobDb) allJobs() ([]*pb.Job, error) {
	q := `SELECT id,
		type,
		content,
		userId,
		destination,
		minute,
		hour,
		dayOfMonth,
		month,
		dayOfWeek,
		year,
		enabled,
		created::INT
	FROM job
	WHERE enabled = true`

	jobs, err := j.queryJobs(q)

	return jobs, err

}

func (j *jobDb) Disable(id string) error {
	q := `UPDATE job SET enabled = false WHERE id = $1`

	_, err := j.Db.Exec(q, id)

	return err
}

func (j *jobDb) queryJobs(q string, dest ...interface{}) ([]*pb.Job, error) {
	jobs := []*pb.Job{}

	rows, err := j.Db.Query(q, dest...)
	if err != nil {
		return jobs, err
	}

	defer rows.Close()
	for rows.Next() {
		job := pb.Job{Schedule: &pb.Schedule{}}
		err := rows.Scan(
			&job.Id,
			&job.Type,
			&job.Content,
			&job.UserId,
			&job.Destination,
			&job.Schedule.Minute,
			&job.Schedule.Hour,
			&job.Schedule.DayOfMonth,
			&job.Schedule.Month,
			&job.Schedule.DayOfWeek,
			&job.Schedule.Year,
			&job.Enabled,
			&job.Created,
		)
		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *jobDb) queryJob(q string, dest ...interface{}) (*pb.Job, error) {
	jobs, err := j.queryJobs(q, dest...)
	if err != nil {
		return nil, err
	}

	if len(jobs) < 1 {
		return nil, nil
	}

	return jobs[0], nil
}
