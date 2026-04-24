package meilisearch

import (
	"context"
	"testing"

	msearch "github.com/meilisearch/meilisearch-go"
)

func TestMeilisearchIndexCreated(t *testing.T) {
	textIndexName := "test-index-name"

	m, err := NewTestMeiliSearch(t)

	if err != nil {
		t.Fatalf("failed to create meilisearhc container: %v", err)
	}

	ctx := context.Background()
	_, err = m.Client.Client.CreateIndexWithContext(ctx, &msearch.IndexConfig{
		Uid:        textIndexName,
		PrimaryKey: "id",
	})

	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
}
