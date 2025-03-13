package controller

type Bot interface {
	SendMessage(rs *Resource, message Message, chatID string, processors ...MessageProcessor) error
	ServeBotAPI(rs *Resource)
}

type Message interface {
	GetContent() string
}
