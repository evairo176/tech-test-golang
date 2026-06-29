package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// Config holds application configuration
type Config struct {
	Port           string
	DBPath         string
	MigrationsPath string
}

// Load returns application config from environment or defaults
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Port:           port,
		DBPath:         "./persons.db",
		MigrationsPath: "file://database/migrations",
	}
}

// InitDB opens a database connection, runs migrations & seeders
func InitDB(cfg *Config) *sql.DB {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	RunMigrations(db, cfg.MigrationsPath)
	RunSeeders(db)

	return db
}

// RunMigrations applies all up migrations
func RunMigrations(db *sql.DB, migrationsPath string) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}

// RunSeeders inserts initial data if the table is empty
func RunSeeders(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Person").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check Person count: %v", err)
	}

	if count > 0 {
		log.Println(fmt.Sprintf("Seeders skipped (%d records already exist)", count))
		return
	}

	_, err = db.Exec(`
		INSERT INTO Person (Name, Country) VALUES
			('Adam', 'Kuala Lumpur'),
			('John', 'Singapore'),
			('Henry', 'Singapore'),
			('Dominic', 'Thailand');
	`)
	if err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Println("Seeders applied successfully")
}
