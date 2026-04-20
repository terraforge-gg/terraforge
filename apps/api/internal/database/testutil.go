package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/terraforge-gg/terraforge/internal/logger"
)

type TestDatabase struct {
	Container   *postgres.PostgresContainer
	Db          *sql.DB
	DatabaseUrl string
	Logger      *slog.Logger
}

func NewTestDatabase(t *testing.T) (*TestDatabase, error) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	databaseUrl, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := NewPostgresConnection(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	logger := logger.New()

	return &TestDatabase{
		Container:   container,
		Db:          db,
		DatabaseUrl: databaseUrl,
		Logger:      logger,
	}, nil
}

const TestUser1Id = "1"
const TestUser1Username = "testuser1"
const TestUser1Email = "testuser1@example.com"

const TestUser2Id = "2"
const TestUser2Username = "testuser2"
const TestUser2Email = "testuser2@example.com"

func (td *TestDatabase) Setup() error {
	if err := Migrate(td.Logger, td.Db); err != nil {
		return err
	}
	ctx := context.Background()
	_, err := InsertTestUser(ctx, td.Db, TestUser1Id, TestUser1Username, TestUser1Email)
	_, err = InsertTestUser(ctx, td.Db, TestUser2Id, TestUser2Username, TestUser2Email)

	return err
}

func (td *TestDatabase) Teardown(ctx context.Context) error {
	if td.Db != nil {
		td.Db.Close()
	}
	return td.Container.Terminate(ctx)
}

func (td *TestDatabase) CleanDatabase(ctx context.Context) error {
	query := `
		DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public')
			LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`
	_, err := td.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}
	return nil
}

func NewTestDatabaseWithCleanup(t *testing.T) (*TestDatabase, error) {
	t.Helper()

	td, err := NewTestDatabase(t)
	if err != nil {
		return nil, err
	}

	if err := td.Setup(); err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		ctx := context.Background()
		if err := td.Teardown(ctx); err != nil {
			t.Logf("failed to teardown test database: %v", err)
		}
	})

	return td, nil
}

func InsertTestUser(ctx context.Context, db *sql.DB, id string, username, email string) (string, error) {
	now := time.Now()

	_, err := db.ExecContext(ctx, `
		INSERT INTO "user" ("id", "name", "username", "email", "emailVerified", "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, username, username, email, true, now, now)
	if err != nil {
		return "", fmt.Errorf("failed to insert test user: %w", err)
	}

	return id, nil
}
