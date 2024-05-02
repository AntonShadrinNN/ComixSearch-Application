// Package stem provides utilities for proccessing data.
package stem

import (
	"comixsearch/internal/models"
	"context"
)

// Normilizer defines a method for normilizing data.
//
//go:generate mockery --name=Normalizer
type Normalizer interface {
	Normalize(ctx context.Context, comices []models.Comic, maxProc int) ([]models.Comic, error)
}
