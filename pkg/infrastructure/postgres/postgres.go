package postgres

import (
	"database/sql"
	"fmt"

	"super-indo-api/pkg/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return db, nil
}

// RunMigration menjalankan SQL migration string ke database
func RunMigration(db *sql.DB, migrationSQL string) error {
	_, err := db.Exec(migrationSQL)
	return err
}
