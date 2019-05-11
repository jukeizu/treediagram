package processor

import (
	"database/sql"
	"fmt"

	"github.com/jukeizu/treediagram/api/protobuf-spec/processingpb"
	"github.com/jukeizu/treediagram/pkg/processor/migrations"
	"github.com/shawntoffel/gossage"
)

const (
	DatabaseName = "treediagram_processor"
)

type Repository interface {
	SaveProcessingRequest(*processingpb.ProcessingRequest) error
	ProcessingRequest(string) (*processingpb.ProcessingRequest, error)
	SaveProcessingEvent(*processingpb.ProcessingEvent) error
	ProcessingEvents(string) ([]*processingpb.ProcessingEvent, error)
	SaveMessageReply(*processingpb.MessageReply) error
	MessageReply(string) (*processingpb.MessageReply, error)

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

func (r *repository) SaveProcessingRequest(p *processingpb.ProcessingRequest) error {
	q := `INSERT INTO processing_request (
		type,
		intentId,
		source,
		channelId,
		serverId,
		botId,
		userId
	) VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id, created::INT`

	err := r.Db.QueryRow(q,
		p.Type,
		p.IntentId,
		p.Source,
		p.ChannelId,
		p.ServerId,
		p.BotId,
		p.UserId,
	).Scan(&p.Id, &p.Created)

	return err
}

func (r *repository) ProcessingRequest(id string) (*processingpb.ProcessingRequest, error) {
	q := `SELECT id,
		type,
		intentId,
		source,
		channelId,
		serverId,
		botId,
		userId,
		created::INT
	FROM processing_request
	WHERE id = $1`

	p := processingpb.ProcessingRequest{}

	err := r.Db.QueryRow(q, id).Scan(
		&p.Id,
		&p.Type,
		&p.IntentId,
		&p.Source,
		&p.ChannelId,
		&p.ServerId,
		&p.BotId,
		&p.UserId,
		&p.Created,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repository) SaveProcessingEvent(e *processingpb.ProcessingEvent) error {
	q := `INSERT INTO processing_event (
		processingRequestId,
		description,
		type
	) VALUES ($1,$2,$3)
	RETURNING id, created::INT`

	err := r.Db.QueryRow(q,
		e.ProcessingRequestId,
		e.Description,
		e.Type,
	).Scan(&e.Id, &e.Created)

	return err
}

func (r *repository) ProcessingEvents(processingRequestId string) ([]*processingpb.ProcessingEvent, error) {
	q := `SELECT id,
		processingRequestId,
		description,
		type,
		created::INT
	FROM processing_event
	WHERE processingRequestId = $1`

	processingEvents := []*processingpb.ProcessingEvent{}

	rows, err := r.Db.Query(q, processingRequestId)
	if err != nil {
		return processingEvents, err
	}

	defer rows.Close()
	for rows.Next() {
		e := processingpb.ProcessingEvent{}

		err := rows.Scan(
			&e.Id,
			&e.ProcessingRequestId,
			&e.Description,
			&e.Type,
			&e.Created,
		)
		if err != nil {
			return processingEvents, err
		}

		processingEvents = append(processingEvents, &e)
	}

	return processingEvents, nil
}

func (r *repository) SaveMessageReply(mr *processingpb.MessageReply) error {
	q := `INSERT INTO message_reply (
		processingRequestId,
		channelId,
		userId,
		content
	) VALUES ($1,$2,$3,$4)
	RETURNING id, created::INT`

	err := r.Db.QueryRow(q,
		mr.ProcessingRequestId,
		mr.ChannelId,
		mr.UserId,
		mr.Content,
	).Scan(&mr.Id, &mr.Created)

	return err
}

func (r *repository) MessageReply(id string) (*processingpb.MessageReply, error) {
	q := `SELECT id,
		processingRequestId,
		channelId,
		userId,
		content,
		created::INT
	FROM message_reply
	WHERE id = $1`

	mr := processingpb.MessageReply{}

	err := r.Db.QueryRow(q, id).Scan(
		&mr.Id,
		&mr.ProcessingRequestId,
		&mr.ChannelId,
		&mr.UserId,
		&mr.Content,
		&mr.Created,
	)

	if err != nil {
		return nil, err
	}

	return &mr, nil
}
