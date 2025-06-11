package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/domains"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddDomains, downAddDomains)
}

func upAddDomains(ctx context.Context, tx *sql.Tx) error {
	postgres.DB_MIGRATOR.CreateTable(&domains.Domain{})
	postgres.DB_MIGRATOR.CreateTable(&domains.SubDomain{})
	return nil
}

func downAddDomains(ctx context.Context, tx *sql.Tx) error {
	postgres.DB_MIGRATOR.DropTable(&domains.SubDomain{})
	postgres.DB_MIGRATOR.DropTable(&domains.Domain{})
	return nil
}
