package dto

import (
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type CreateProjectRequest struct {
	Name    string  `json:"name" validate:"required,min=3,max=100"`
	Slug    string  `json:"slug" validate:"required,min=3,max=100,url_slug"`
	Summary *string `json:"summary" validate:"omitempty,max=120"`
	Type    string  `json:"type" validate:"project_type"`
}

type ProjectResponse struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Summary     *string   `json:"summary"`
	Description *string   `json:"description"`
	IconUrl     *string   `json:"iconUrl"`
	Downloads   int64     `json:"downloads"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserId      string    `json:"userId"`
}

type ProjectMemberResponse struct {
	Username string  `json:"username"`
	UserId   string  `json:"userId"`
	Image    *string `json:"image"`
	Role     string  `json:"role"`
}

func ProjectToProjectResponse(p models.Project) ProjectResponse {
	return ProjectResponse{
		Id:          p.Id,
		Name:        p.Name,
		Slug:        p.Slug,
		Summary:     p.Summary,
		Description: p.Description,
		IconUrl:     p.IconUrl,
		Downloads:   p.Downloads,
		Type:        string(p.Type),
		Status:      string(p.Status),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		UserId:      p.UserId,
	}
}

func ProjectMemberToProjectMemberResponse(member models.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		UserId:   member.UserId,
		Role:     string(member.Role),
		Username: member.User.Username,
		Image:    member.User.Image,
	}
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,url_slug,min=3,max=100"`
	Summary     *string `json:"summary,omitempty" validate:"omitempty,max=120"`
	Description *string `json:"description,omitempty"`
	IconUrl     *string `json:"iconUrl,omitempty" validate:"omitempty,url"`
}
