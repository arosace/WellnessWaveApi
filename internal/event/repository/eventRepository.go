package repository

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
)

type EventRepository interface {
	Add(ctx context.Context, event model.Event) (*model.Event, error)
	GetByHealthSpecialistId(context.Context, string) ([]*model.Event, error)
	GetByPatientId(context.Context, string) ([]*model.Event, error)
	Update(context.Context, model.RescheduleRequest) (*model.Event, error)
}
