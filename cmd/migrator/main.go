package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		log.Fatal("storage-path is required")
	}
	if migrationsPath == "" {
		log.Fatal("migrations-path is required")
	}

	// Создаем экземпляр драйвера для подключения к PostgreSQL
	driver, err := postgres.WithInstance(openDB(storagePath), &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create database instance: %v", err)
	}

	// Создаем новый мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}

	// Применяем миграции
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatalf("failed to apply migrations: %v", err)
	}

	fmt.Println("migrations applied successfully")
}

// openDB открывает соединение с базой данных
func openDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil
	}
	return db
}
