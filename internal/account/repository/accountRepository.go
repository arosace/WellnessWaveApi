package repository

import (
	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/labstack/echo/v5"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx echo.Context, user model.Account) (*model.Account, error)
	Update(ctx echo.Context, user *model.Account) (*model.Account, error)
	UpdateAuth(ctx echo.Context, user *model.Account) (*model.Account, error)
	List(ctx echo.Context) ([]*model.Account, error)
	FindByID(ctx echo.Context, id string) (*model.Account, error)
	FindByEmail(ctx echo.Context, email string) (*model.Account, error)
	FindByParentID(ctx echo.Context, parentId string) ([]*model.Account, error)
}
