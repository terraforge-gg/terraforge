package models

import "time"

type User struct {
	Id              string
	Name            string
	Username        string
	DisplayUsername *string
	Email           string
	EmailVerified   bool
	Image           *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
