package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/servers"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDropUsernameServers, downDropUsernameServers)
}

func upDropUsernameServers(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropColumn(&servers.Server{}, "username")
}

func downDropUsernameServers(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN username VARCHAR(255)")
	return err
}
