package models

import "gorm.io/gorm"

const RECEIVERTYPESINGLE = "single"
const RECEIVERTYPEALL = "all"

type Subscriber struct {
	gorm.Model
	ChatID   string `gorm:"uniqueIndex"`
	UserName string
	Status   string
	Channel  string
}

type IncomingMessage struct {
	gorm.Model
	ChatID         string
	UserName       string
	MessageText    string
	MessageCommand string
}

// TODO: Remove GeneralMessage field from the model, as it's not relavant
type OutgoingMessage struct {
	gorm.Model
	ReceiverChatID string
	AyahID         *uint
	GeneralMessage string
}

type Category struct {
	gorm.Model
	CategoryEnglish string
	CategoryBangla  string
}

type VersePreference struct {
	gorm.Model
	VerseId int
	Status  string
}
