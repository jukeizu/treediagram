package migration

import (
	"database/sql"
)

type CreateTableProcessingEvent20190122002825 struct{}

func (m CreateTableProcessingEvent20190122002825) Version() string {
	return "20190122002825_CreateTableProcessingEvent"
}

func (m CreateTableProcessingEvent20190122002825) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS processing_event (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			processingRequestId UUID NOT NULL REFERENCES processing_request,
			description STRING NOT NULL DEFAULT '',
			type STRING NOT NULL DEFAULT '',
			timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	return err
}

func (m CreateTableProcessingEvent20190122002825) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE processing_event`)
	return err
}
