package internal

import (
	"comixsearch/internal/models"
	"comixsearch/internal/normalizer"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/html"
)

const (
	UrlArchive = "https://xkcd.com/archive/"
	UrlComix   = "https://xkcd.com/%s/info.0.json"
	// ComixCount = 2882
)

type Servicer interface {
	Run() error
	Stop() error
	Read() error
	Write() error
	ProcessData([]models.Comix) error
	GetData() ([]models.Comix, error)
	Search(http.ResponseWriter, *http.Request)
}

type ServiceXKCD struct {
}

func (s *ServiceXKCD) Search(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world\n")
}

func (s *ServiceXKCD) Run() {
	http.HandleFunc("/", s.ProcessData)
	http.ListenAndServe(":8080", nil)
}

func (s *ServiceXKCD) ProcessData(w http.ResponseWriter, r *http.Request) {
	comixes, _ := s.GetData()
	for _, comix := range comixes {
		title, _ := normalizer.Normalize(comix.Title)
		content, _ := normalizer.Normalize(comix.Content)
		altContent, _ := normalizer.Normalize(comix.AltContent)
		c := models.DBUnit{
			Id:         comix.Id,
			Title:      title,
			Content:    content,
			AltContent: altContent,
			Link:       comix.Link,
		}

		file, err := os.OpenFile("data.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Ошибка открытия файла:", err)
			return
		}
		defer file.Close()

		// Создаем JSON-кодировщик
		encoder := json.NewEncoder(file)

		// Записываем данные в файл
		if err := encoder.Encode(c); err != nil {
			fmt.Println("Ошибка записи в файл:", err)
			return
		}

		fmt.Println("Данные успешно записаны в файл data.json")
	}
}

func (s *ServiceXKCD) GetData() ([]models.Comix, error) {

	resp, err := http.Get(UrlArchive)
	if err != nil {
		fmt.Println("Ошибка при загрузке страницы:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Парсим HTML документ
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при парсинге HTML:", err)
		return nil, err
	}

	// Ищем все теги <a> внутри <div id="middleContainer">
	links := findLinksInDiv(doc, "middleContainer")
	comixes := make([]models.Comix, 0)
	for i, link := range links {
		if i == 1 {
			return comixes, nil
		}

		url := fmt.Sprintf(UrlComix, link)
		r, err := http.Get(url)

		if err != nil {
			fmt.Println(err)
		}

		defer r.Body.Close()

		var c models.Comix
		io.Copy(os.Stdout, r.Body)
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&c)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		comixes = append(comixes, c)
	}

	return comixes, err
}

func findLinksInDiv(n *html.Node, id string) []string {
	var links []string

	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					re := regexp.MustCompile(`/(\d+)/$`)
					link := re.FindString(attr.Val)
					re = regexp.MustCompile(`/`)
					link = re.ReplaceAllString(link, "")

					links = append(links, link)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	var findDiv func(*html.Node)
	findDiv = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					findLinks(n)
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findDiv(c)
		}
	}

	findDiv(n)

	return links
}
