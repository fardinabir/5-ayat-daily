package controller

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/db/repos"
	"one-minute-quran/models"
)

func init() {
	LoadFromConfig()
}

type Resource struct {
	Bot   Bot
	Store *repos.Store
}

func NewResource() *Resource {
	tb, err := NewTgBot()
	if err != nil {
		return nil
	}
	ss := repos.NewSubsStore()
	rs := &Resource{
		Bot:   tb,
		Store: ss,
	}
	return rs
}

func (rs *Resource) PublishToSubscribers(ayah *models.Ayah) error {
	ayahText := FormatAyahText(ayah)

	subscribersList, err := rs.Store.GetAllSubscribers()
	for _, subscriber := range subscribersList {
		err = rs.Bot.SendMessage(rs, ayahText, subscriber.ChatID, &ayah.ID)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func FormatAyahText(ayah *models.Ayah) string {
	return fmt.Sprintf("%v,\n(%v),\n(%v)\n[%v:%v]", ayah.AyahTextArabic, ayah.AyahTextBangla,
		ayah.AyahTextEnglish, ayah.SuraNo, ayah.VerseNo)
}

func (rs *Resource) ServeBot() {
	rs.Bot.ServeBotAPI(rs)
}

func LoadFromConfig() {
	viper.SetConfigFile("./config/.config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
