package processor

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processing"
	migration "github.com/jukeizu/treediagram/processor/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_processor"
)

type Repository interface {
	SaveProcessingRequest(*processing.ProcessingRequest) error
	ProcessingRequest(string) (*processing.ProcessingRequest, error)
	SaveProcessingEvent(*processing.ProcessingEvent) error
	ProcessingEvents(string) ([]*processing.ProcessingEvent, error)
	SaveMessageReply(*processing.MessageReply) error
	MessageReply(string) (*processing.MessageReply, error)

	Migrate() error
}

type repository struct {
	Db *sql.DB
}

func NewRepository(url string) (Repository, error) {
	conn := fmt.Sprintf("postgresql://%s/%s?sslmode=disable", url, DatabaseName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	r := repository{
		Db: db,
	}

	return &r, nil
}

func (r *repository) Migrate() error {
	_, err := r.Db.Exec(`CREATE DATABASE IF NOT EXISTS ` + DatabaseName)
	if err != nil {
		return err
	}

	g, err := gossage.New(r.Db)
	if err != nil {
		return err
	}

	err = g.RegisterMigrations(
		migration.CreateTableProcessingRequest20190121234940{},
		migration.CreateTableProcessingEvent20190122002825{},
		migration.CreateTableMessageReply20190122003737{},
	)
	if err != nil {
		return err
	}

	return g.Up()
}

func (r *repository) SaveProcessingRequest(p *processing.ProcessingRequest) error {
	return nil
}

func (r *repository) ProcessingRequest(id string) (*processing.ProcessingRequest, error) {
	return &processing.ProcessingRequest{}, nil
}

func (r *repository) SaveProcessingEvent(e *processing.ProcessingEvent) error {
	return nil
}

func (r *repository) ProcessingEvents(requestId string) ([]*processing.ProcessingEvent, error) {
	return []*processing.ProcessingEvent{}, nil
}

func (r *repository) SaveMessageReply(mr *processing.MessageReply) error {
	return nil
}

func (r *repository) MessageReply(id string) (*processing.MessageReply, error) {
	return nil, nil
}
