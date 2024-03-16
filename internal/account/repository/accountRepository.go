package repository

import (
	"context"

	"../model"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx context.Context, user model.Account) error
	FindByID(ctx context.Context, id string) (*model.Account, error)
}
