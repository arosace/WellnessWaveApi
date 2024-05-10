package repository

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/labstack/echo/v5"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx echo.Context, user model.Account) (*model.Account, error)
	Update(ctx context.Context, user *model.Account) (*model.Account, error)
	UpdateAuth(ctx context.Context, user *model.Account) (*model.Account, error)
	List(ctx echo.Context) ([]*model.Account, error)
	FindByID(ctx context.Context, id string) (*model.Account, error)
	FindByEmail(ctx echo.Context, email string) (*model.Account, error)
	FindByParentID(ctx context.Context, parentId string) ([]*model.Account, error)
}
