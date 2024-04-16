package interfaces

import (
	"comixsearch/internal/models"
	"context"
)

type Storager interface {
	Get(ctx context.Context, keywords []string, isContentSearch bool) ([]string, error)
	Write(ctx context.Context, data models.Comic) error
	Close()
}
