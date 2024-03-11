package controller

import (
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/controller/interfaces"
	"one-minute-quran/db"
	"one-minute-quran/db/repos"
	"one-minute-quran/models"
)

func init() {
	LoadFromConfig()
}

const SubscriberChanBuffer int = 10

type Resource struct {
	Bot       interfaces.Bot
	SubsStore *repos.SubscriberStore
}

func NewResource() *Resource {
	tgBot := models.NewTgBot()
	ss := &repos.SubscriberStore{DB: db.ConnectDB()}
	return &Resource{
		Bot:       tgBot,
		SubsStore: ss,
	}
}

func (rs *Resource) PublishToSubscribers(message string) error {
	subscribersList, err := rs.SubsStore.GetAllSubscribers()
	for _, subscriber := range subscribersList {
		err = rs.Bot.SendMessage(message, subscriber.ChatID)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (rs *Resource) ServeBot() {
	newSubscriber := make(chan models.Subscriber, SubscriberChanBuffer)
	go rs.Bot.ServeBot(newSubscriber)

	for subscriber := range newSubscriber {
		err := rs.SubsStore.Save(&subscriber)
		if err != nil {
			log.Println("Cannot Save Subscriber Info: ", err)
		} else {
			log.Println("New subscriber added successfully")
		}
	}
}

func LoadFromConfig() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
