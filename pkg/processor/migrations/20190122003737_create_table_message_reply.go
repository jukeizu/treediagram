package migrations

import (
	"database/sql"
)

type CreateTableMessageReply20190122003737 struct{}

func (m CreateTableMessageReply20190122003737) Version() string {
	return "20190122003737_CreateTableMessageReply"
}

func (m CreateTableMessageReply20190122003737) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS message_reply (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			processingRequestId UUID NOT NULL REFERENCES processing_request,
			channelId STRING NOT NULL DEFAULT '',
			userId STRING NOT NULL DEFAULT '',
			content JSONB,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	return err
}

func (m CreateTableMessageReply20190122003737) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE message_reply`)
	return err
}
