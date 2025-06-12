package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/deployments"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateDeployments, downCreateDeployments)
}

func upCreateDeployments(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.CreateTable(&deployments.Deployment{})
}

func downCreateDeployments(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropTable(&deployments.Deployment{})
}
