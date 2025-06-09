package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddSecretName, downAddSecretName)
}

func upAddSecretName(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE secrets ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT ''")
	if err != nil {
		return err
	}
	return nil
}

func downAddSecretName(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE secrets DROP COLUMN IF EXISTS name")
	if err != nil {
		return err
	}
	return nil
}
