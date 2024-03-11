package interfaces

import "one-minute-quran/models"

type Bot interface {
	SendMessage(message, chatID string) error
	ServeBot(chan<- models.Subscriber)
}
