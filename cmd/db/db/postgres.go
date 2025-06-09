package postgres

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"deployer.com/modules/secrets"
	"deployer.com/modules/users"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GORM_DB *gorm.DB
var SQL_DB *sql.DB
var DB_MIGRATOR gorm.Migrator
var DNS = os.Getenv("DATABASE_URL")

func NewGormDBMigration() error {
	db, err := gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err == nil {
		GORM_DB = db
		SQL_DB, _ = db.DB()
		DB_MIGRATOR = db.Migrator()
	}
	return err
}

func NewGormDB(lc fx.Lifecycle) (*gorm.DB, error) {
	// dsn := os.Getenv("DATABASE_URL")
	if DNS == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	GORM_DB = db
	SQL_DB, err = db.DB()
	if err != nil {
		return nil, err
	}
	DB_MIGRATOR = db.Migrator()

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			dbSQL, err := db.DB()
			if err == nil {
				return dbSQL.Close()
			}
			return nil
		},
	})
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&users.User{}, &secrets.Secret{})
}
