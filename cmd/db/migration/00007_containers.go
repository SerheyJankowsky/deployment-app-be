package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/containers"
)

func init() {
	goose.AddMigrationContext(upContainers, downContainers)
}

func upContainers(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.CreateTable(&containers.Container{})
}

func downContainers(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropTable(&containers.Container{})
}
