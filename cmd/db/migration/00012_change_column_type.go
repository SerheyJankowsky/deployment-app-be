package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upChangeColumnType, downChangeColumnType)
}

func upChangeColumnType(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE scripts ALTER COLUMN last_run_at DROP NOT NULL, ALTER COLUMN last_run_at SET DEFAULT NULL`)
	if err != nil {
		return err
	}
	return nil
}

func downChangeColumnType(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE scripts ALTER COLUMN last_run_at SET NOT NULL, ALTER COLUMN last_run_at DROP DEFAULT`)
	if err != nil {
		return err
	}
	return nil
}
