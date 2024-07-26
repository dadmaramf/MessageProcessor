package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "migrations to storage")
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of migrations table")

	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	if storagePath == "" {
		panic("storage-path is required")
	}

	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("postgres://%s?sslmode=disable", storagePath))

	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")

}
