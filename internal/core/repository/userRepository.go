package repository

import (
	"context"

	"../model"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Add(ctx context.Context, user model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
}
