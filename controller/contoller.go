package controller

import (
	"five-ayat-daily/db/repos"
	"five-ayat-daily/models"
	"fmt"
	"github.com/spf13/viper"
	"log"
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

func (rs *Resource) PublishToSubscribers(ayah *models.Ayah, generalMsg string) error {
	ayahText := FormatAyahText(ayah)

	subscribersList, err := rs.Store.GetAllSubscribers()
	log.Println("fetched total subscriber : ", len(subscribersList))
	for _, subscriber := range subscribersList {
		if generalMsg != "" {
			err = rs.Bot.SendMessage(rs, generalMsg, subscriber.ChatID, nil)
		} else {
			err = rs.Bot.SendMessage(rs, ayahText, subscriber.ChatID, &ayah.ID)
		}
		if err != nil {
			log.Println(fmt.Sprintf("error while sending msg to : %v, chatID : %v", subscriber.UserName, subscriber.ChatID), err)
		}
	}
	return nil
}

func FormatAyahText(ayah *models.Ayah) string {
	if ayah == nil {
		return ""
	}
	return fmt.Sprintf("%v,\n(%v),\n(%v)\n[%v:%v]", ayah.AyahTextArabic, ayah.AyahTextBangla,
		ayah.AyahTextEnglish, ayah.SuraNo, ayah.VerseNo)
}

func (rs *Resource) ServeBot() {
	log.Println("serving bot....")
	rs.Bot.ServeBotAPI(rs)
}

func LoadFromConfig() {
	viper.SetConfigFile("./config/.config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
