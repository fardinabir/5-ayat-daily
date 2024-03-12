package controller

type Bot interface {
	SendMessage(rs *Resource, message, chatID string, ayahId *uint) error
	ServeBotAPI(rs *Resource)
}
