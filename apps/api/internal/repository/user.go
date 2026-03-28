package repository

import (
	"context"
	"database/sql"

	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type UserRepository interface {
	FindUserByIdentifier(ctx context.Context, q database.Querier, userIdentifier string) (*models.User, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindUserByIdentifier(ctx context.Context, q database.Querier, userIdentifier string) (*models.User, error) {
	query := `
		SELECT
			"id",
			"name",
			"username",
			"displayUsername",
			"email",
			"emailVerified",
			"image",
			"createdAt",
			"updatedAt"
		FROM "user" 
		WHERE ("id" = $1 OR "username" = $1);`

	user := &models.User{}
	err := q.QueryRowContext(ctx, query, userIdentifier).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.DisplayUsername,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
