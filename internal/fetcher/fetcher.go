package fetcher

import (
	"comixsearch/internal/models"
	"context"
)

// Fetcher defines a method for retrieving comic data with specified parameters.
type Fetcher interface {
	GetData(ctx context.Context, maxProc int, lastId int64) ([]models.Comic, error)
}
