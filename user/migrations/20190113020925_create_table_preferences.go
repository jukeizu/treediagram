package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190113020925, Down20190113020925)
}

func Up20190113020925(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS preferences (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			userId STRING NOT NULL DEFAULT '',
			serverId STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW())`)
	return err
}

func Down20190113020925(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE preferences`)
	return err
}