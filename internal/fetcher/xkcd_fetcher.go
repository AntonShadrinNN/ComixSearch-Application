// Package fetcher is implementing a data fetching functionality for retrieving comic information from the XKCD website.
package fetcher

import (
	"comixsearch/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// XkcdFetcher is used to fetch XKCD comic information from specified URLs using an HTTP client.
type XkcdFetcher struct {
	UrlArchive string      //  for accessing the XKCD comic archive.
	UrlComic   string      // is the URL used to fetch the latest XKCD comic from the XKCD website
	client     http.Client // is used to make HTTP requests when fetching data from the XKCD API.
}

// The NewFetcher function initializes an XkcdFetcher instance.
func NewFetcher(UrlArchive, UrlComix string) XkcdFetcher {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	return XkcdFetcher{
		UrlArchive: UrlArchive,
		UrlComic:   UrlComix,
		client:     c,
	}
}

// GetData fetches comic data from the XKCD website.
func (p XkcdFetcher) GetData(ctx context.Context, maxProc int, lastId int64) ([]models.Comic, error) {
	newLastId, err := p.getComicesCount()
	if err != nil {
		return nil, err
	}

	comicesCnt := newLastId - lastId
	wp := NewWorkerPool(maxProc, int(comicesCnt))
	for id := lastId + 1; id <= newLastId; id++ {
		wp.submit(id)
	}
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(wp.start(p.UrlComic))

	comices := make([]models.Comic, 0, comicesCnt)

	eg.Go(
		func() error {
			for i := 0; i < int(comicesCnt); i++ {
				res := wp.getResult()
				if res.err != nil {
					log.Printf("Enable to get comic. Error: %s", res.err)
					continue
				}
				comices = append(comices, res.data)
			}

			return nil
		},
	)

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return comices, nil
}

// getComicesCount fetches the total number of comics available in the XKCD archive.
func (p XkcdFetcher) getComicesCount() (int64, error) {
	req, err := http.NewRequest("GET", p.UrlArchive, nil)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}

	var ar archiveInfo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ar)
	if err != nil {
		return 0, err
	}

	return ar.ComicesCount, nil
}

type archiveInfo struct {
	ComicesCount int64 `json:"num"`
}

type worker struct {
	urlSrc     string
	taskChan   <-chan int64
	resultChan chan<- result
}

func (w *worker) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for id := range w.taskChan {
			data, err := fetch(w.urlSrc, id)
			w.resultChan <- result{data: data, err: err}
		}
	}()
}

type result struct {
	data models.Comic
	err  error
}

// Fetch retrieves a comic from a specified URL using the provided ID.
func fetch(urlSrc string, id int64) (models.Comic, error) {
	url := fmt.Sprintf(urlSrc, id)
	log.Println(url)
	r, err := http.Get(url)

	if err != nil {
		return models.Comic{}, err
	}

	defer r.Body.Close()

	var c models.Comic
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&c)

	if err != nil {
		return models.Comic{}, err
	}
	return c, nil
}

// workerPool is a pool of workers used for concurent fetching comices from comic' website.
type workerPool struct {
	taskChan     chan int64
	resultChan   chan result
	workersCount int
	waitGroup    *sync.WaitGroup
}

func NewWorkerPool(workersCount, tasksCount int) *workerPool {
	return &workerPool{
		taskChan:     make(chan int64, tasksCount),
		resultChan:   make(chan result),
		workersCount: workersCount,
		waitGroup:    &sync.WaitGroup{},
	}
}

func (wp *workerPool) start(urlSrc string) func() error {
	return func() error {
		for i := 0; i < wp.workersCount; i++ {
			worker := worker{
				urlSrc:     urlSrc,
				taskChan:   wp.taskChan,
				resultChan: wp.resultChan,
			}
			worker.start(wp.waitGroup)
		}

		close(wp.taskChan)
		defer close(wp.resultChan)
		wp.waitGroup.Wait()
		return nil
	}
}

func (wp *workerPool) submit(id int64) {
	wp.taskChan <- id
}
func (wp *workerPool) getResult() result {
	return <-wp.resultChan
}
