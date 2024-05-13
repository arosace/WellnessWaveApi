package repository

import (
	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/labstack/echo/v5"
)

type EventRepository interface {
	Add(ctx echo.Context, event model.Event) (*model.Event, error)
	GetByHealthSpecialistId(echo.Context, string, string) ([]*model.Event, error)
	GetByPatientId(echo.Context, string, string) ([]*model.Event, error)
	Update(echo.Context, model.RescheduleRequest) (*model.Event, error)
}
