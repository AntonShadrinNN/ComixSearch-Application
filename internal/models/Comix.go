package models

type Comix struct {
	Id         int    `json:"num"`
	Title      string `json:"title"`
	Content    string `json:"transcript"`
	AltContent string `json:"alt"`
	Link       string `json:"img"`
}
