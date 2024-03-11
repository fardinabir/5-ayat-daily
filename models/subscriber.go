package models

import "gorm.io/gorm"

type Subscriber struct {
	gorm.Model
	ChatID  string `gorm:"uniqueIndex"`
	Status  string
	Channel string
}

type IncomingMessage struct {
	gorm.Model
	ChatID         string
	MessageText    string
	MessageCommand string
}

type OutgoingMessage struct {
	gorm.Model
	MessageText string
	AyahID      uint
}

type Ayah struct {
	gorm.Model
	SuraNo          int
	VerseNo         int
	AyahTextArabic  string
	AyahTextBangla  string
	AyahTextEnglish string
	CategoryID      int
}

type Category struct {
	gorm.Model
	CategoryEnglish string
	CategoryBangla  string
}
