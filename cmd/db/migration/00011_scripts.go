package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/scripts"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upScripts, downScripts)
}

func upScripts(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.CreateTable(&scripts.Script{})
}

func downScripts(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropTable(&scripts.Script{})
}
