package migrations

import (
	"database/sql"
)

type CreateTableJob20190119052738 struct{}

func (m CreateTableJob20190119052738) Version() string {
	return "20190119052738_CreateTableJob"
}

func (m CreateTableJob20190119052738) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS job(
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			userId STRING NOT NULL DEFAULT '',
			source STRING NOT NULL DEFAULT '',
			content STRING NOT NULL DEFAULT '',
			endpoint STRING NOT NULL DEFAULT '',
			destination STRING NOT NULL DEFAULT '',
			minute STRING NOT NULL DEFAULT '',
			hour STRING NOT NULL DEFAULT '',
			dayOfMonth STRING NOT NULL DEFAULT '',
			month STRING NOT NULL DEFAULT '',
			dayOfWeek STRING NOT NULL DEFAULT '',
			year STRING NOT NULL DEFAULT '',
			enabled BOOL NOT NULL DEFAULT TRUE,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)
	return err
}

func (m CreateTableJob20190119052738) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE job`)
	return err
}
