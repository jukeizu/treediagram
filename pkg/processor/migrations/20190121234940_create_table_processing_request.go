package migrations

import (
	"database/sql"
)

type CreateTableProcessingRequest20190121234940 struct{}

func (m CreateTableProcessingRequest20190121234940) Version() string {
	return "20190121234940_CreateTableProcessingRequest"
}

func (m CreateTableProcessingRequest20190121234940) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS processing_request (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			type STRING NOT NULL DEFAULT '',
			intentId STRING NOT NULL DEFAULT '',
			source STRING NOT NULL DEFAULT '',
			channelId STRING NOT NULL DEFAULT '',
			serverId STRING NOT NULL DEFAULT '',
			botId STRING NOT NULL DEFAULT '',
			userId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	return err
}

func (m CreateTableProcessingRequest20190121234940) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE processing_request`)
	return err
}
