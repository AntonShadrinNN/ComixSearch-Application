package models

type Comic struct {
	Id      int64  `json:"num"`
	Title   string `json:"title"`
	Content string `json:"transcript"`
	// AltContent string `json:"alt"`
	Link string `json:"img"`
}
