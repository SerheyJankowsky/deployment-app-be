package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddUsernameServers, downAddUsernameServers)
}

func upAddUsernameServers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN IF NOT EXISTS username VARCHAR(255) NOT NULL DEFAULT ''")
	if err != nil {
		return err
	}
	return nil
}

func downAddUsernameServers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers DROP COLUMN IF EXISTS username")
	if err != nil {
		return err
	}
	return nil
}
