package controller

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/controller/interfaces"
	"one-minute-quran/db/repos"
	"one-minute-quran/models"
)

func init() {
	LoadFromConfig()
}

type Resource struct {
	Bot       interfaces.Bot
	SubsStore *repos.SubscriberStore
}

func NewResource() *Resource {
	tgBot := models.NewTgBot()
	ss := repos.NewSubsStore()
	rs := &Resource{
		Bot:       tgBot,
		SubsStore: ss,
	}
	tgBot.Rs = rs
	return rs
}

func (rs *Resource) PublishToSubscribers(ayah *models.Ayah) error {
	ayahText := FormatAyahText(ayah)

	rs.SubsStore.SaveOutgoingMessage(&models.OutgoingMessage{
		ReceiverType: models.RECEIVERTYPEALL,
		AyahID:       ayah.ID,
	})

	subscribersList, err := rs.SubsStore.GetAllSubscribers()
	for _, subscriber := range subscribersList {
		err = rs.Bot.SendMessage(ayahText, subscriber.ChatID)
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
	go rs.Bot.ServeBot()
}

func LoadFromConfig() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
