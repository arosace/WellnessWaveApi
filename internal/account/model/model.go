package model

import "context"

type Model interface {
	ValidateModel(ctx context.Context) error
}
