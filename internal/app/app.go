package internal

import (
	"comixsearch/internal/app/interfaces"
	"comixsearch/internal/models"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// const (
// 	UrlArchive = "https://xkcd.com/archive/"
// 	UrlComix   = "https://xkcd.com/%s/info.0.json"
// 	// ComicsCount = 2882
// )

type SearchApp struct {
	N          interfaces.Normalizer
	P          interfaces.Fetcher
	S          interfaces.Storager
	L          interfaces.Logger
	maxProc    int
	UrlArchive string
	UrlComix   string
}

func NewApp(N interfaces.Normalizer, P interfaces.Fetcher, S interfaces.Storager, L interfaces.Logger, maxProc int) SearchApp {
	return SearchApp{
		N:       N,
		P:       P,
		S:       S,
		L:       L,
		maxProc: maxProc,
	}
}

// func (p XkcdParser) ProcessData() {
// 	comixes, _ := getData()
// 	for _, comix := range comixes {
// 		title, _ := p.N.Normalize(comix.Title)
// 		content, _ := p.N.Normalize(comix.Content)
// 		altContent, _ := p.N.Normalize(comix.AltContent)
// 		c := models.DBUnit{
// 			Id:         comix.Id,
// 			Title:      title,
// 			Content:    content,
// 			AltContent: altContent,
// 			Link:       comix.Link,
// 		}

// 		file, err := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
// 		if err != nil {
// 			fmt.Println("Ошибка открытия файла:", err)
// 			return
// 		}
// 		defer file.Close()

// 		// Создаем JSON-кодировщик
// 		encoder := json.NewEncoder(file)

// 		// Записываем данные в файл
// 		if err := encoder.Encode(c); err != nil {
// 			fmt.Println("Ошибка записи в файл:", err)
// 			return
// 		}

// 		fmt.Println("Данные успешно записаны в файл data.json")
// 	}
// }

// func (s SearchApp) updateDB() error {
// 	ch := make(chan models.Comic)

// 	go func() {
// 		s.P.GetData(ch, 5)
// 		close(ch)
// 	}()
// 	fmt.Println("Get end")
// 	c := 0
// 	for m := range ch {
// 		fmt.Println(m)
// 		// _, err := s.N.Normalize(m.Title)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// }
// 		c++
// 		// title, err := s.N.Normalize(m.Title)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return err
// 		// }
// 		// content, _ := s.N.Normalize(m.Content)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return err
// 		// }
// 		// altContent, _ := s.N.Normalize(m.AltContent)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return err
// 		// }
// 		// err = s.S.Write(m.Id, title, content, altContent, m.Link)
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// 	return err
// 		// }

// 	}

// 	// some log messages here
// 	// fmt.Println("End func")
// 	return nil
// }

func (s SearchApp) ProcessData(ctx context.Context, in <-chan *models.Comic, out chan<- *models.Comic) error {
	eg, _ := errgroup.WithContext(ctx)

	// sem := make(chan struct{}, s.maxProc)

	// for comic := range in {
	// 	if comic.Id == 0 {
	// 		out <- &models.Comic{} // Сигнализируем об окончании передачи комиксов
	// 		break
	// 	}

	// 	// Блокируем выполнение, если количество запущенных горутин равно максимальному
	// 	sem <- struct{}{}
	// 	eg.Go(
	// 		func() error {
	// 			defer func() { <-sem }()
	// 			var (
	// 				stemmedData string
	// 				err         error
	// 			)
	// 			stemmedData, err = s.N.Normalize(comic.Content)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			comic.Content = stemmedData
	// 			stemmedData, err = s.N.Normalize(comic.Title)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			comic.Title = stemmedData
	// 			out <- comic
	// 			return nil
	// 		},
	// 	)
	// }

	// // End:
	// if err := eg.Wait(); err != nil {
	// 	return err
	// }
	// return nil

	workersCount := 0
	for {
		for workersCount <= s.maxProc {
			comic := <-in
			if comic.Id == 0 {
				out <- &models.Comic{}
				goto End
			}
			eg.Go(
				func() error {
					defer func() { workersCount-- }()
					var (
						stemmedData string
						err         error
					)
					stemmedData, err = s.N.Normalize(comic.Content)
					if err != nil {
						return err
					}

					comic.Content = stemmedData
					stemmedData, err = s.N.Normalize(comic.Title)
					if err != nil {
						return err
					}

					comic.Title = stemmedData
					out <- comic
					return nil
				},
			)
		}
	}
End:
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s SearchApp) writeToDatabase(ctx context.Context, ch <-chan *models.Comic) error {
	eg, _ := errgroup.WithContext(ctx)

	workersCount := 0
	for {
		for workersCount <= s.maxProc {
			workersCount++
			comic := <-ch

			if comic.Id == 0 {
				goto End
			}

			eg.Go(
				func() error {
					defer func() { workersCount-- }()
					err := s.S.Write(ctx, *comic)
					if err != nil {
						return err
					}
					return nil
				},
			)
		}
	}
End:
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s SearchApp) Search(ctx context.Context, keywords []string, isContent bool) ([]string, error) {
	foundLinks, err := s.S.Get(ctx, keywords, isContent)
	return foundLinks, err
}

func (s SearchApp) Run(ctx context.Context) func() error {
	// go func() {
	// 	for {
	// 		// FIXME: implement isUpdates()
	// 		// if s.isUpdates() {
	// 		// FIXME: Replace log with our custom Logger
	// 		log.Fatal(s.updateDB())
	// 		time.Sleep(time.Hour)
	// 		// }
	// 	}
	// }()
	return func() error {
		fmt.Println("Start")
		eg, _ := errgroup.WithContext(ctx)
		comics := make(chan *models.Comic, 15)

		processedComics := make(chan *models.Comic, 15)
		defer close(processedComics)
		defer close(comics)

		eg.Go(
			func() error {
				err := s.P.GetData(ctx, comics, s.maxProc)
				if err != nil {
					return err
				}
				return nil
			},
		)

		eg.Go(
			func() error {
				err := s.ProcessData(ctx, comics, processedComics)
				if err != nil {
					return err
				}

				return nil
			},
		)

		eg.Go(
			func() error {
				err := s.writeToDatabase(ctx, processedComics)
				if err != nil {
					return err
				}

				return nil
				// cnt := 0
				// for comic := range processedComics {
				// 	fmt.Println(comic)
				// 	if comic.Id == 0 {
				// 		return nil
				// 	}
				// 	// cnt++
				// }
				// return nil
			},
		)

		if err := eg.Wait(); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("End")
		return nil
		// http.HandleFunc("/", s.search)
		// http.ListenAndServe(":8080", nil)
	}
}
