package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/servers"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddServers, downAddServers)
}

func upAddServers(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.CreateTable(&servers.Server{})
}

func downAddServers(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropTable(&servers.Server{})
}
