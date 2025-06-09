package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/secrets"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSecrets, downCreateSecrets)
}

func upCreateSecrets(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return postgres.DB_MIGRATOR.CreateTable(&secrets.Secret{})
}

func downCreateSecrets(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return postgres.DB_MIGRATOR.DropTable(&secrets.Secret{})
}
