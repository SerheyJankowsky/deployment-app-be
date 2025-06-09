package main

import (
	"log"

	postgres "deployer.com/cmd/db/db"
	_ "deployer.com/cmd/db/migration"
	"github.com/pressly/goose/v3"
)

func main() {
	err := postgres.NewGormDBMigration()
	if err != nil {
		log.Fatal(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(postgres.SQL_DB, "cmd/db/migration"); err != nil {
		panic(err)
	}
}
