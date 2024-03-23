package repository

import (
	"context"

	"github.com/arosace/WellnessWaveApi/internal/event/model"
)

type EventRepository interface {
	Add(ctx context.Context, event model.Event) (*model.Event, error)
}
