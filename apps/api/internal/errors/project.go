package errors

import "errors"

var (
	ErrProjectUnauthorisedAction = errors.New("unauthorised project action")
	ErrProjectNotFound           = errors.New("project not found")
	ErrProjectSlugUsed           = errors.New("project slug is not available")
)
