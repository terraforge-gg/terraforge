package meilisearch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestMeiliSearch struct {
	Container testcontainers.Container
	Client    *MeiliSearchClient
	HostUrl   string
	MasterKey string
}

func NewTestMeiliSearch(t *testing.T) (*TestMeiliSearch, error) {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "getmeili/meilisearch:latest",
		ExposedPorts: []string{"7700/tcp"},
		Env: map[string]string{
			"MEILI_MASTER_KEY": "masterKey",
		},
		WaitingFor: wait.ForHTTP("/health").WithPort("7700/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start meilisearch container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "7700/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	hostUrl := fmt.Sprintf("http://%s:%s", host, port.Port())
	masterKey := "masterKey"

	client := NewMeiliSearch(hostUrl, masterKey)

	return &TestMeiliSearch{
		Container: container,
		Client:    client,
		HostUrl:   hostUrl,
		MasterKey: masterKey,
	}, nil
}

func (tm *TestMeiliSearch) Teardown(ctx context.Context) error {
	return tm.Container.Terminate(ctx)
}

func NewTestMeiliSearchWithCleanup(t *testing.T) (*TestMeiliSearch, error) {
	t.Helper()

	tm, err := NewTestMeiliSearch(t)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		ctx := context.Background()
		if err := tm.Teardown(ctx); err != nil {
			t.Logf("failed to teardown test meilisearch: %v", err)
		}
	})

	return tm, nil
}
