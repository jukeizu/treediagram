package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190113072028, Down20190113072028)
}

func Up20190113072028(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS intents(
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

func Down20190113072028(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE intents`)
	return err
}
