package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUdpateSecrets, downUdpateSecrets)
}

func upUdpateSecrets(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE secrets ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}
	_, err = tx.Exec("ALTER TABLE secrets ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	return err
}

func downUdpateSecrets(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE secrets DROP COLUMN IF EXISTS updated_at")
	if err != nil {
		return err
	}
	_, err = tx.Exec("ALTER TABLE secrets DROP COLUMN IF EXISTS deleted_at")
	return err
}
