package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/servers"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRemoveServersKey, downRemoveServersKey)
}

func upRemoveServersKey(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropColumn(&servers.Server{}, "key")
}

func downRemoveServersKey(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN key VARCHAR(255)")
	return err
}
