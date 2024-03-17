package repository

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx context.Context, user model.Account) error
	List(ctx context.Context) ([]model.Account, error)
	FindByID(ctx context.Context, id string) (*model.Account, error)
	FindByEmail(ctx context.Context, email string) (*model.Account, error)
}
