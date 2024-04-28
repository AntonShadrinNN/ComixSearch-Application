package stem

import (
	"comixsearch/internal/models"
	"context"
	"log"
	"strings"
	"sync"

	"github.com/kljensen/snowball"
	"golang.org/x/sync/errgroup"
)

// Stem is used for stemming data.
type Stem struct {
	language  string // which language is used for stemming.
	stopWords bool   // used or not stop words.
}

// The NewStem function initializes an Stem instance.
func NewStem(language string, stopWords bool) Stem {
	return Stem{
		language:  language,
		stopWords: stopWords,
	}
}

// Proccess stems preproccessed data and join result in string.
func (s Stem) process(data string) (string, error) {
	tokens := prepocess(data)

	for i := 0; i < len(tokens); i++ {
		stemmed, err := snowball.Stem(tokens[i], s.language, s.stopWords)
		tokens[i] = stemmed
		if err != nil {
			return "", err
		}
	}

	stemmedData := strings.Join(tokens, " ")

	return stemmedData, nil
}

// Normilize gets comices and stems each of them.
func (s Stem) Normalize(ctx context.Context, comices []models.Comic, maxProc int) ([]models.Comic, error) {
	wp := NewWorkerPool(maxProc, len(comices))
	for _, comic := range comices {
		wp.submit(comic)
	}

	eg, _ := errgroup.WithContext(ctx)
	eg.Go(wp.start(s))

	eg.Go(
		func() error {
			for i := 0; i < int(len(comices)); i++ {
				res := wp.GetResult()
				if res.err != nil {
					log.Printf("Enable to get comic. Error: %s", res.err)
					return res.err
				}
				comices[i] = res.data
			}

			return nil
		},
	)

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return comices, nil
}

type worker struct {
	s          Stem
	taskChan   <-chan models.Comic
	resultChan chan<- result
}

func (w *worker) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for comic := range w.taskChan {
			title, err := w.s.process(comic.Title)
			if err != nil {
				w.resultChan <- result{
					models.Comic{},
					err,
				}
				continue
			}
			content, err := w.s.process(comic.Content)
			temp := models.Comic{
				Id:      comic.Id,
				Title:   title,
				Content: content,
				Link:    comic.Link,
			}
			w.resultChan <- result{data: temp, err: err}
		}
	}()
}

type result struct {
	data models.Comic
	err  error
}

// workerPool is a pool of workers used for concurent normilizing data.
type workerPool struct {
	taskChan     chan models.Comic
	resultChan   chan result
	workersCount int
	waitGroup    *sync.WaitGroup
}

func NewWorkerPool(workersCount, tasksCount int) *workerPool {
	return &workerPool{
		taskChan:     make(chan models.Comic, tasksCount),
		resultChan:   make(chan result),
		workersCount: workersCount,
		waitGroup:    &sync.WaitGroup{},
	}
}

func (wp *workerPool) start(s Stem) func() error {
	return func() error {
		for i := 0; i < wp.workersCount; i++ {
			worker := worker{
				s:          s,
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

func (wp *workerPool) submit(comic models.Comic) {
	wp.taskChan <- comic
}
func (wp *workerPool) GetResult() result {
	return <-wp.resultChan
}
