package migration

import (
	"database/sql"
)

type CreateTableIntent20190113072028 struct{}

func (m CreateTableIntent20190113072028) Version() string {
	return "20190113072028_CreateTableIntent"
}

func (m CreateTableIntent20190113072028) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS intent(
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			serverId STRING NOT NULL DEFAULT '',
			name STRING NOT NULL DEFAULT '',
			regex STRING NOT NULL DEFAULT '',
			mention BOOL NOT NULL DEFAULT FALSE,
			response STRING NOT NULL DEFAULT '',
			endpoint STRING NOT NULL DEFAULT '',
			help STRING NOT NULL DEFAULT '',
			enabled BOOL NOT NULL DEFAULT TRUE,
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)
	return err
}

func (m CreateTableIntent20190113072028) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE intent`)
	return err
}
