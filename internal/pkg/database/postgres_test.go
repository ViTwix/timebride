package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"timebride/internal/config"
)

func TestNewPostgresConnection(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "postgres"
	cfg.Database.Password = "secret"
	cfg.Database.DBName = "timebride"

	db, err := NewPostgresConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.HealthCheck(); err != nil {
		t.Errorf("Database health check failed: %v", err)
	}
}

func TestDatabase_Transaction(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.Database.User = "postgres"
	cfg.Database.Password = "secret"
	cfg.Database.DBName = "timebride"

	db, err := NewPostgresConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Тест успішної транзакції
	err = db.Transaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "SELECT 1")
		return err
	})
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}

	// Тест відкату транзакції
	err = db.Transaction(ctx, func(tx *sql.Tx) error {
		return fmt.Errorf("test error")
	})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
