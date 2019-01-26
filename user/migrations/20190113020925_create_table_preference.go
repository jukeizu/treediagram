package migration

import (
	"database/sql"
)

type CreateTablePreference20190113020925 struct{}

func (m CreateTablePreference20190113020925) Version() string {
	return "20190113020925_CreateTablePreference"
}

func (m CreateTablePreference20190113020925) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS preference (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			userId STRING UNIQUE NOT NULL DEFAULT '',
			serverId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)
	return err
}

func (m CreateTablePreference20190113020925) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE preference`)
	return err
}
