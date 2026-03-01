package meilisearch

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
	msearch "github.com/meilisearch/meilisearch-go"
)

type MeiliSearchClient struct {
	Client msearch.ServiceManager
}

func NewMeiliSearch(meiliSeachHostUrl string, apiKey string) *MeiliSearchClient {
	return &MeiliSearchClient{
		Client: meilisearch.New(meiliSeachHostUrl, meilisearch.WithAPIKey(apiKey)),
	}
}

func (m *MeiliSearchClient) Health(ctx context.Context) error {
	_, err := m.Client.HealthWithContext(ctx)

	if err != nil {
		return err
	}

	return nil
}
