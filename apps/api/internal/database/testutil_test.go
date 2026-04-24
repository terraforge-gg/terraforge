package database_test

import (
	"context"
	"testing"

	"github.com/terraforge-gg/terraforge/internal/database"
)

func TestTestDatabaseSetup(t *testing.T) {
	td, err := database.NewTestDatabaseWithCleanup(t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := td.Db.Ping(); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Log("successfully connected to test database")
}

func TestTestDatabaseCleanup(t *testing.T) {
	td, err := database.NewTestDatabaseWithCleanup(t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	_, err = td.Db.Exec("CREATE TABLE test_cleanup (id SERIAL PRIMARY KEY)")
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	ctx := context.Background()
	if err := td.CleanDatabase(ctx); err != nil {
		t.Fatalf("failed to clean database: %v", err)
	}

	t.Log("successfully cleaned test database")
}

func TestTestDatabaseSetupInsertsTestUser(t *testing.T) {
	td, err := database.NewTestDatabaseWithCleanup(t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	ctx := context.Background()
	var count int
	err = td.Db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "user"`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query user count: %v", err)
	}

	if count != 2 {
		t.Fatalf("expected 2 users, got %d", count)
	}

	var email string
	err = td.Db.QueryRowContext(ctx, `SELECT "email" FROM "user"`).Scan(&email)
	if err != nil {
		t.Fatalf("failed to query user email: %v", err)
	}

	if email != database.TestUser1Email {
		t.Fatalf("expected email testuser1@example.com, got %s", email)
	}

	t.Log("successfully verified test user was inserted")
}

func TestTestDatabaseSetupInsertsTestLoaderVersion(t *testing.T) {
	td, err := database.NewTestDatabaseWithCleanup(t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	ctx := context.Background()
	var count int
	err = td.Db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "loader_version"`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query loader version count: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 loader version, got %d", count)
	}

	var gameVersion string
	err = td.Db.QueryRowContext(ctx, `SELECT "gameVersion" FROM "loader_version"`).Scan(&gameVersion)
	if err != nil {
		t.Fatalf("failed to query loader version game version: %v", err)
	}

	if gameVersion != database.TestLoaderVersionGameVersion {
		t.Fatalf("expected game veresion 1.4.4, got %s", gameVersion)
	}

	t.Log("successfully verified test loader version was inserted")
}
