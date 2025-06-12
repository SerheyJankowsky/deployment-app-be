package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upServersRemoveCertificate, downServersRemoveCertificate)
}

func upServersRemoveCertificate(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM information_schema.columns
				WHERE table_name='servers' AND column_name='certificate'
			) THEN
				ALTER TABLE servers DROP COLUMN certificate;
			END IF;
		END
		$$;
	`)
	return err
}

func downServersRemoveCertificate(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN certificate VARCHAR(255)")
	return err
}
