package errors

import "errors"

var (
	ErrProjectUnauthorisedAction = errors.New("unauthorized project action")
	ErrProjectNotFound           = errors.New("project not found")
	ErrProjectSlugUsed           = errors.New("project slug is not available")
	ErrNoProjectsFound           = errors.New("no projects found")
)
