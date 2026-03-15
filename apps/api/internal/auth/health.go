package auth

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type AuthHealthCheckService struct {
	logger     *slog.Logger
	BaseUrl    string
	httpClient *http.Client
}

func NewAuthHealthCheckService(logger *slog.Logger, baseUrl string) *AuthHealthCheckService {
	return &AuthHealthCheckService{
		BaseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *AuthHealthCheckService) Health(ctx context.Context) error {
	resp, err := s.httpClient.Head(s.BaseUrl + "/health")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
