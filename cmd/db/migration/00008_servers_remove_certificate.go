package migrations

import (
	"context"
	"database/sql"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/servers"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upServersRemoveCertificate, downServersRemoveCertificate)
}

func upServersRemoveCertificate(ctx context.Context, tx *sql.Tx) error {
	return postgres.DB_MIGRATOR.DropColumn(&servers.Server{}, "certificate")
}

func downServersRemoveCertificate(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("ALTER TABLE servers ADD COLUMN certificate VARCHAR(255)")
	return err
}
