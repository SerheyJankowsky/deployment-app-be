package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/users"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUser, downCreateUser)
}

func upCreateUser(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return postgres.DB_MIGRATOR.CreateTable(&users.User{})
}

func downCreateUser(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return postgres.DB_MIGRATOR.DropTable(&users.User{})
}
