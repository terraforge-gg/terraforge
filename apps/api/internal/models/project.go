package models

import (
	"time"
)

type ProjectType string

const ProjectTypeMod ProjectType = "mod"

type ProjectStatus string

const (
	ProjectStatusDraft    ProjectStatus = "draft"
	ProjectStatusRejected ProjectStatus = "rejected"
	ProjectStatusApproved ProjectStatus = "approved"
	ProjectStatusBanned   ProjectStatus = "banned"
)

type Project struct {
	Id          string
	Name        string
	Slug        string
	Summary     *string
	Description *string
	IconUrl     *string
	Downloads   int64
	Type        ProjectType
	Status      ProjectStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	UserId      string
}

type ProjectMemberRole string

const (
	ProjectMemberRoleOwner      ProjectMemberRole = "owner"
	ProjectMemberRoleAdmin      ProjectMemberRole = "admin"
	ProjectMemberRoleDeveloper  ProjectMemberRole = "developer"
	ProjectMemberRoleMaintainer ProjectMemberRole = "maintainer"
	ProjectMemberRoleMember     ProjectMemberRole = "member"
)

type ProjectMember struct {
	Id        string
	ProjectId string
	UserId    string
	Role      ProjectMemberRole
	CreatedAt time.Time
	User      User
}
