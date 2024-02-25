package interfaces

type Bot interface {
	SendMessage(message, chatID string) error
	ServeBot() interface{}
}
