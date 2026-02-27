package database

import "errors"

var (
	ErrUniqueViolation = errors.New("unique violation")
)
