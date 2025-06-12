package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upIndexFix, downIndexFix)
}

func upIndexFix(ctx context.Context, tx *sql.Tx) error {
	// Удаляем лишний уникальный индекс на name (если он мешает)
	_, err := tx.ExecContext(ctx, `
		DROP INDEX IF EXISTS idx_sub_domains_name
	`)
	if err != nil {
		return err
	}

	// Добавляем обычный индекс на domain_id (если его еще нет)
	// Используем IF NOT EXISTS для безопасности
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_sub_domains_domain_id ON sub_domains (domain_id)
	`)
	return err
}

func downIndexFix(ctx context.Context, tx *sql.Tx) error {
	// При откате удаляем индекс на domain_id
	_, err := tx.ExecContext(ctx, `
		DROP INDEX IF EXISTS idx_sub_domains_domain_id
	`)
	if err != nil {
		return err
	}

	// Восстанавливаем обычный (не уникальный) индекс на name
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_sub_domains_name ON sub_domains (name)
	`)
	return err
}
