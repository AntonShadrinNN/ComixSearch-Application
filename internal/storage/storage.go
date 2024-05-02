// Package storage specifies a set of methods that a storage implementation must have.
package storage

import (
	"comixsearch/internal/models"
	"context"
)

// Storager is used for interacting with a storage system in a consistent way.
//
//go:generate mockery --name=Storager
type Storager interface {
	GetComices(ctx context.Context, keywords []string, limit int) (map[string]string, error)
	Write(ctx context.Context, data []models.Comic) error
	GetLastId(ctx context.Context) (int64, error)
	Close()
}
