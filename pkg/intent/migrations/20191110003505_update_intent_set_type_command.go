package migrations

import "database/sql"

type UpdateIntentSetTypeCommand20191110003505 struct{}

func (m UpdateIntentSetTypeCommand20191110003505) Version() string {
	return "20191110003505_UpdateIntentSetTypeCommand"
}

func (m UpdateIntentSetTypeCommand20191110003505) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE intent SET type = 'command'`)
	return err
}

func (m UpdateIntentSetTypeCommand20191110003505) Down(tx *sql.Tx) error {
	return nil
}
