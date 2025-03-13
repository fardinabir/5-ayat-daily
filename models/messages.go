package models

import (
	"fmt"
	"gorm.io/gorm"
)

type GeneralMessage struct {
	Message string
}

func (g GeneralMessage) GetContent() string {
	return g.Message
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

func (a Ayah) GetContent() string {
	if a.AyahTextArabic == "" {
		return ""
	}
	return fmt.Sprintf("%v,\n(%v),\n(%v)\n[%v:%v]", a.AyahTextArabic, a.AyahTextBangla,
		a.AyahTextEnglish, a.SuraNo, a.VerseNo)
}
