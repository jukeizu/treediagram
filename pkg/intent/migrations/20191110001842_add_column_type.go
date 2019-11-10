package migrations

import "database/sql"

type AddColumnType20191110001842 struct{}

func (m AddColumnType20191110001842) Version() string {
	return "20191110001842_AddColumnType"
}

func (m AddColumnType20191110001842) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE intent ADD COLUMN type STRING NOT NULL DEFAULT ''`)
	return err
}

func (m AddColumnType20191110001842) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE intent DROP COLUMN type`)
	return err
}
