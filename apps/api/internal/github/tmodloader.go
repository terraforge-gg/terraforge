package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	tModLoaderReleasesURL = "https://api.github.com/repos/tModLoader/tModLoader/releases"
	releasesPerPage       = 100
)

type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Body        string    `json:"body"`
}

func FetchTModLoaderReleases(ctx context.Context) ([]GitHubRelease, error) {
	var allReleases []GitHubRelease
	page := 1

	client := &http.Client{Timeout: 30 * time.Second}

	for {
		releases, hasMore, err := fetchReleasesPage(ctx, client, page)
		if err != nil {
			return nil, err
		}

		allReleases = append(allReleases, releases...)

		if !hasMore {
			break
		}

		page++
	}

	return allReleases, nil
}

func fetchReleasesPage(ctx context.Context, client *http.Client, page int) ([]GitHubRelease, bool, error) {
	params := url.Values{}
	params.Set("per_page", strconv.Itoa(releasesPerPage))
	params.Set("page", strconv.Itoa(page))

	reqURL := fmt.Sprintf("%s?%s", tModLoaderReleasesURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "terraforge-gg")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("GitHub API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, false, fmt.Errorf("failed to decode response: %w", err)
	}

	hasMore := len(releases) == releasesPerPage

	return releases, hasMore, nil
}
