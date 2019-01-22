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
	SaveProcessingRequest(processing.ProcessingRequest) (*processing.ProcessingRequest, error)
	ProcessingRequest(string) (*processing.ProcessingRequest, error)
	SaveProcessingEvent(*processing.ProcessingEvent) error
	ProcessingEvents(string) ([]*processing.ProcessingEvent, error)
	SaveMessageReply(processing.MessageReply) (*processing.MessageReply, error)
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

func (r *repository) SaveProcessingRequest(p processing.ProcessingRequest) (*processing.ProcessingRequest, error) {
	q := `INSERT INTO processing_request (
		intentId,
		source,
		channelId,
		serverId,
		botId,
		userId
	) VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id, created::INT`

	err := r.Db.QueryRow(q,
		p.IntentId,
		p.Source,
		p.ChannelId,
		p.ServerId,
		p.BotId,
		p.UserId,
	).Scan(&p.Id, &p.Created)

	if err != nil {
		return nil, err
	}

	return &p, err
}

func (r *repository) ProcessingRequest(id string) (*processing.ProcessingRequest, error) {
	q := `SELECT id,
			intentId,
			source,
			channelId,
			serverId,
			botId,
			userId,
			created::INT
		FROM processing_request
		WHERE id = $1`

	p := processing.ProcessingRequest{}

	err := r.Db.QueryRow(q, id).Scan(
		&p.Id,
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

func (r *repository) SaveProcessingEvent(e *processing.ProcessingEvent) error {
	q := `INSERT INTO processing_event (
			processingRequestId,
			description,
			type
		) VALUES ($1,$2,$3)`

	_, err := r.Db.Exec(q,
		e.ProcessingRequestId,
		e.Description,
		e.Type,
	)

	return err
}

func (r *repository) ProcessingEvents(processingRequestId string) ([]*processing.ProcessingEvent, error) {
	q := `SELECT id,
			processingRequestId,
			description,
			type,
			timestamp::INT
		FROM processing_event
		WHERE processingRequestId = $1`

	processingEvents := []*processing.ProcessingEvent{}

	rows, err := r.Db.Query(q, processingRequestId)
	if err != nil {
		return processingEvents, err
	}

	defer rows.Close()
	for rows.Next() {
		e := processing.ProcessingEvent{}

		err := rows.Scan(
			&e.Id,
			&e.ProcessingRequestId,
			&e.Description,
			&e.Type,
			&e.Timestamp,
		)
		if err != nil {
			return processingEvents, err
		}

		processingEvents = append(processingEvents, &e)
	}

	return processingEvents, nil
}

func (r *repository) SaveMessageReply(mr processing.MessageReply) (*processing.MessageReply, error) {
	q := `INSERT INTO message_reply (
			processingRequestId,
			channelId,
			userId,
			isPrivateMessage,
			isRedirect,
			content
		) VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created::INT`

	err := r.Db.QueryRow(q,
		mr.ProcessingRequestId,
		mr.ChannelId,
		mr.UserId,
		mr.IsPrivateMessage,
		mr.IsRedirect,
		mr.Content,
	).Scan(&mr.Id, &mr.Created)

	if err != nil {
		return nil, err
	}

	return &mr, err
}

func (r *repository) MessageReply(id string) (*processing.MessageReply, error) {
	q := `SELECT id,
			processingRequestId,
			channelId,
			userId,
			isPrivateMessage,
			isRedirect,
			content,
			created::INT
		FROM message_reply
		WHERE id = $1`

	mr := processing.MessageReply{}

	err := r.Db.QueryRow(q, id).Scan(
		&mr.Id,
		&mr.ProcessingRequestId,
		&mr.ChannelId,
		&mr.UserId,
		&mr.IsPrivateMessage,
		&mr.IsRedirect,
		&mr.Content,
		&mr.Created,
	)

	if err != nil {
		return nil, err
	}

	return &mr, nil
}
