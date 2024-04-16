package internal

import (
	"comixsearch/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

const (
	UrlArchive = "https://xkcd.com/archive/"
	UrlComix   = "https://xkcd.com/%d/info.0.json"
)

type XkcdFetcher struct {
	LastId int64
}

func (p *XkcdFetcher) GetData(ctx context.Context, ch chan<- *models.Comic, maxProc int) error {
	resp, err := http.Get(UrlArchive)
	if err != nil {
		fmt.Println("Ошибка при загрузке страницы:", err)
		return err
	}
	defer resp.Body.Close()
	// Парсим HTML документ
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при парсинге HTML:", err)
		return err
	}

	ids := make(chan int64)
	defer close(ids)
	// Ищем все теги <a> внутри <div id="middleContainer">
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(findLinksInDiv(doc, "middleContainer", ids, p.LastId))
	workersCount := 0
	for {
		for workersCount < maxProc {
			workersCount++
			id := <-ids
			if id > p.LastId {
				p.LastId = id
			}
			// found last new comix
			if id == 0 {
				// signal "End input"

				goto End
			}
			eg.Go(func() error {
				defer func() { workersCount-- }()

				url := fmt.Sprintf(UrlComix, id)
				r, err := http.Get(url)

				if err != nil {
					return err
				}

				defer r.Body.Close()

				var c models.Comic
				dec := json.NewDecoder(r.Body)
				err = dec.Decode(&c)

				if err != nil {
					return err
				}
				ch <- &c
				return nil
			},
			)

		}
	}

	// sem := make(chan struct{}, maxProc)

	// for id := range ids {
	// 	if id == 0 {
	// 		ch <- &models.Comic{} // Сигнализируем об окончании передачи комиксов
	// 		break
	// 	}

	// 	// Блокируем выполнение, если количество запущенных горутин равно максимальному
	// 	sem <- struct{}{}
	// 	url := fmt.Sprintf(UrlComix, id)
	// 	eg.Go(func() error {
	// 		defer func() { <-sem }() // Освобождаем слот в канале после завершения горутины
	// 		r, err := http.Get(url)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer r.Body.Close()

	// 		var c models.Comic
	// 		dec := json.NewDecoder(r.Body)
	// 		if err := dec.Decode(&c); err != nil {
	// 			return err
	// 		}
	// 		ch <- &c // Отправляем комикс в канал
	// 		return nil
	// 	})
	// }

End:
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func findLinksInDiv(n *html.Node, id string, ids chan<- int64, lastInd int64) func() error {
	return func() error {
		endSearch := false
		var findLinks func(*html.Node) error
		findLinks = func(n *html.Node) error {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						re := regexp.MustCompile(`/(\d+)/$`)
						id := re.FindString(attr.Val)
						re = regexp.MustCompile(`/`)
						id = re.ReplaceAllString(id, "")
						intId, err := strconv.ParseInt(id, 10, 64)
						// default value to exit from getData loop
						if intId <= lastInd && lastInd != 0 {
							ids <- 0
							endSearch = true
							return nil
						} else if err == nil {
							ids <- intId
						} else {
							return err
						}

					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := findLinks(c); err != nil {
					return err
				}

				if endSearch {
					return nil
				}
			}
			return nil
		}

		// var findDiv func(*html.Node) error
		var findDiv func(n *html.Node) error
		findDiv = func(n *html.Node) error {
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, attr := range n.Attr {
					if attr.Key == "id" && attr.Val == id {
						if err := findLinks(n); err != nil {
							return err
						}
						return nil
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := findDiv(c); err != nil {
					return err
				}
			}
			return nil
		}
		return findDiv(n)
	}
}
