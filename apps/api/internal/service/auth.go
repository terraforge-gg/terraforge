package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/terraforge-gg/terraforge/internal/auth"
)

type AuthService struct {
	logger     *slog.Logger
	BaseUrl    string
	httpClient *http.Client
}

func NewAuthService(logger *slog.Logger, baseUrl string) *AuthService {
	return &AuthService{
		BaseUrl: baseUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *AuthService) Health(ctx context.Context) error {
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

func (s *AuthService) ValidateSession(cookie *http.Cookie) (*auth.Session, error) {
	req, err := http.NewRequest("GET", s.BaseUrl+"/api/auth/get-session", nil)

	if err != nil {
		s.logger.Error("Failed to create request to auth service", "error", err)
		return nil, err
	}

	req.AddCookie(cookie)

	res, err := s.httpClient.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		s.logger.Error("Failed to validate session with auth service", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	var sessionDto auth.Session
	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		return nil, readErr
	}

	if string(body) == "null" {
		return nil, errors.New("session response is null")
	}

	if err := json.Unmarshal(body, &sessionDto); err != nil {
		return nil, err
	}

	session := auth.Session{
		User: auth.SessionUser{
			Id:          sessionDto.User.Id,
			Name:        sessionDto.User.Name,
			Email:       sessionDto.User.Email,
			Image:       sessionDto.User.Image,
			CreatedAt:   sessionDto.User.CreatedAt,
			UpdatedAt:   sessionDto.User.UpdatedAt,
			Username:    sessionDto.User.Username,
			DisplayName: sessionDto.User.DisplayName,
		},
	}

	return &session, nil
}
