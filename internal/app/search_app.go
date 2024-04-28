// Package app contains
// the implementation of a search application for fetching, processing, and storing comic data.
package app

import (
	"comixsearch/internal/app/interfaces"
	"comixsearch/internal/config"
	"comixsearch/internal/fetcher"
	"comixsearch/internal/models"
	"comixsearch/internal/stem"
	"comixsearch/internal/storage"
	"context"
	"log"
)

// SearchApp encapsulates components for normalizing, fetching, storing data, logging, and
// setting maximum processing capacity.
type SearchApp struct {
	normalizer stem.Normalizer   // is likely used for text normalization.
	fetcher    fetcher.Fetcher   // is responsible for retrieving data from external sources.
	storage    storage.Storager  // is responsible for storing and managing data within the search application.
	logger     interfaces.Logger // is likely used for logging messages, errors, and other information related to the application's operations.
	maxProc    int               // us maximum number of concurrent processes that the application can handle.
}

// The NewApp function initializes an SearchApp instance.
func NewApp(ctx context.Context, conf config.Config) (SearchApp, error) {
	fetcher := fetcher.NewFetcher(conf.UrlArchive, conf.UrlComic)
	stemmer := stem.NewStem("english", true)
	var (
		err    error
		lastId int64
	)

	storage, err := storage.NewStorage(ctx, conf.DbConn)
	if err != nil {
		return SearchApp{}, err
	}
	log.Println("Create storage is success")

	s := SearchApp{
		normalizer: stemmer,
		fetcher:    fetcher,
		storage:    storage,
		maxProc:    conf.MaxProc,
	}

	lastId, err = s.storage.GetLastId(ctx)
	if err != nil {
		return SearchApp{}, err
	}

	comices, err := s.fetchComices(ctx, lastId)
	if err != nil {
		return SearchApp{}, err
	}
	log.Println("Fetch is success")

	comices, err = s.processData(ctx, comices)
	if err != nil {
		return SearchApp{}, err
	}
	log.Println("Process is success")

	err = s.writeToDatabase(ctx, comices)

	log.Println("Write to database is success")
	return s, err
}

// fetchComices fetches comic data using the fetcher component.
func (s SearchApp) fetchComices(ctx context.Context, lastId int64) ([]models.Comic, error) {
	return s.fetcher.GetData(ctx, s.maxProc, lastId)
}

// ProcessData normalizes the comic data before further processing.
func (s SearchApp) processData(ctx context.Context, comices []models.Comic) ([]models.Comic, error) {
	return s.normalizer.Normalize(ctx, comices, s.maxProc)
}

// writeToDatabase write the comic data to the database using the `storage` component.
func (s SearchApp) writeToDatabase(ctx context.Context, comices []models.Comic) error {
	return s.storage.Write(ctx, comices)

}

// Search allows searching for comic data based on provided keywords and a limit.
func (s SearchApp) Search(ctx context.Context, keywords []string, limit int) (map[string]string, error) {
	return s.storage.GetComices(ctx, keywords, limit)
}
