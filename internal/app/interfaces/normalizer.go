package interfaces

// const (
// 	Language  = "english"
// 	StopWords = true
// )

type Normalizer interface {
	Normalize(string) (string, error)
}
