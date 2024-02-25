package models

import "gorm.io/gorm"

type Subscriber struct {
	gorm.Model
	ChatID  string `gorm:"uniqueIndex"`
	Status  string
	Channel string
}
