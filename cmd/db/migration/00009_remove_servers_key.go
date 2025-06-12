package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRemoveServersKey, downRemoveServersKey)
}

func upRemoveServersKey(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name='servers' AND column_name='key'
			) THEN
				ALTER TABLE servers DROP COLUMN key;
			END IF;
		END
		$$;
	`)
	return err
}

func downRemoveServersKey(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN key VARCHAR(255)")
	return err
}
