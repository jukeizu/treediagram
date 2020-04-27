package migrations

import "database/sql"

type AlterJobAddColumnInstanceid20200427030530 struct{}

func (m AlterJobAddColumnInstanceid20200427030530) Version() string {
	return "20200427030530_AlterJobAddColumnInstanceid"
}

func (m AlterJobAddColumnInstanceid20200427030530) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE job ADD COLUMN instanceId STRING UNIQUE NOT NULL DEFAULT from_uuid(uuid_v4())`)
	return err
}

func (m AlterJobAddColumnInstanceid20200427030530) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE job DROP COLUMN instanceId`)
	return err
}
