package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	redis_client_wrapper "github.com/terraforge-gg/terraforge/internal/lib/redis"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type ProjectCache interface {
	GetProject(ctx context.Context, identifier string) (*models.Project, error)
	SetProject(ctx context.Context, project *models.Project, ttl time.Duration) error
	DeleteProject(ctx context.Context, id string) error
	GetProjectMembers(ctx context.Context, identifier string) ([]models.ProjectMember, error)
	SetProjectMembers(ctx context.Context, project *models.Project, projectMembers []models.ProjectMember, ttl time.Duration) error
}

type cache struct {
	Wrapper *redis_client_wrapper.RedisClient
}

func NewProjectCache(redisWrapper *redis_client_wrapper.RedisClient) ProjectCache {
	return &cache{
		Wrapper: redisWrapper,
	}
}

func (c *cache) GetProject(ctx context.Context, identifier string) (*models.Project, error) {
	// First try the identifier directly as an ID key
	if project, err := c.getProjectById(ctx, identifier); err == nil {
		return project, nil
	}

	// Otherwise treat it as a slug — resolve to an ID first
	id, err := c.Wrapper.Client.Get(ctx, slugKey(identifier)).Result()

	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, fmt.Errorf("redis slug lookup: %w", err)
	}

	p, err := c.getProjectById(ctx, id)

	if errors.Is(err, ErrCacheMiss) {
		return nil, ErrCacheMiss
	}

	return p, err
}

func (c *cache) GetProjectMembers(ctx context.Context, identifier string) ([]models.ProjectMember, error) {
	// First try the identifier directly as an ID key
	if project, err := c.getProjectMembersByProjectId(ctx, identifier); err == nil {
		return project, nil
	}

	// Otherwise treat it as a slug — resolve to an ID first
	id, err := c.Wrapper.Client.Get(ctx, projectMembersSlugKey(identifier)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis slug lookup: %w", err)
	}

	pm, err := c.getProjectMembersByProjectId(ctx, id)

	if errors.Is(err, ErrCacheMiss) {
		return nil, ErrCacheMiss
	}

	return pm, err
}

func (c *cache) SetProject(ctx context.Context, project *models.Project, ttl time.Duration) error {
	b, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("marshal project: %w", err)
	}

	pipe := c.Wrapper.Client.Pipeline()
	// Single source of truth — project data lives under its ID
	pipe.Set(ctx, projectKey(project.Id), b, ttl)
	// Slug is just a pointer to the ID
	pipe.Set(ctx, slugKey(project.Slug), project.Id, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *cache) getProjectById(ctx context.Context, id string) (*models.Project, error) {
	val, err := c.Wrapper.Client.Get(ctx, projectKey(id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var project models.Project
	if err := json.Unmarshal([]byte(val), &project); err != nil {
		return nil, fmt.Errorf("unmarshal project: %w", err)
	}
	return &project, nil
}

func (c *cache) getProjectMembersByProjectId(ctx context.Context, id string) ([]models.ProjectMember, error) {
	val, err := c.Wrapper.Client.Get(ctx, projectMembersKey(id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var projectMembers []models.ProjectMember
	if err := json.Unmarshal([]byte(val), &projectMembers); err != nil {
		return nil, fmt.Errorf("unmarshal project: %w", err)
	}
	return projectMembers, nil
}

func (c *cache) SetProjectMembers(ctx context.Context, project *models.Project, projectMembers []models.ProjectMember, ttl time.Duration) error {
	b, err := json.Marshal(projectMembers)
	if err != nil {
		return fmt.Errorf("marshal project members: %w", err)
	}

	pipe := c.Wrapper.Client.Pipeline()
	// Single source of truth — project data lives under its ID
	pipe.Set(ctx, projectMembersKey(project.Id), b, ttl)
	// Slug is just a pointer to the ID
	pipe.Set(ctx, projectMembersSlugKey(project.Slug), project.Id, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *cache) DeleteProject(ctx context.Context, id string) error {
	if err := c.Wrapper.Client.Del(ctx, projectKey(id)).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}

	return nil
}

func projectKey(id string) string {
	return "project:" + id
}

func projectMembersKey(id string) string {
	return "project:" + id + ":members"
}

func projectMembersSlugKey(slug string) string {
	return "project:slug:" + slug + ":members"
}

func slugKey(slug string) string {
	return "project:slug:" + slug
}
