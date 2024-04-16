package interfaces

import (
	"comixsearch/internal/models"
	"context"
)

type Fetcher interface {
	GetData(ctx context.Context, comics chan<- *models.Comic, maxProc int) error
}
