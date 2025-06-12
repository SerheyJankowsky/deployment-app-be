package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddNameIndex, downAddNameIndex)
}

func upAddNameIndex(ctx context.Context, tx *sql.Tx) error {
	// Проверяем, существует ли индекс, и создаем только если не существует
	_, err := tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_domain_subdomain 
		ON sub_domains (name, domain_id)
	`)
	return err
}

func downAddNameIndex(ctx context.Context, tx *sql.Tx) error {
	// Удаляем индекс при откате миграции
	_, err := tx.ExecContext(ctx, `
		DROP INDEX IF EXISTS idx_domain_subdomain
	`)
	return err
}
