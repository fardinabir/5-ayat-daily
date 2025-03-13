package controller

import (
	"five-ayat-daily/models"
	"log"
)

type MessageProcessor func() error

func WithMessageLogging(chatID string, message Message) MessageProcessor {
	return func() error {
		log.Printf("Message logged in for chatID: %v, message: %v", chatID, message)
		return nil
	}
}

func WithDBPersistence(rs *Resource, om *models.OutgoingMessage) MessageProcessor {
	return func() error {
		return rs.Store.SaveOutgoingMessage(om)
	}
}
