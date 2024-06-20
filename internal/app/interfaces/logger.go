package interfaces

type Logger interface {
	Info(text string)
	Debug(text string)
	Error(text string)
}
