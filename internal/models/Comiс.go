// Package models contains data structures related to the application's domain model.
package models

// The Comic struct defines the structure of a comic with fields for ID, title, content, and link.
type Comic struct {
	Id      int64  `json:"num"`
	Title   string `json:"title"`
	Content string `json:"transcript"`
	Link    string `json:"img"`
}
