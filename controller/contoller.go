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

func (rs *Resource) PublishToSubscribers(msg Message) error {
	subscribersList, err := rs.Store.GetAllSubscribers()
	log.Println("fetched total subscriber : ", len(subscribersList))

	for _, subscriber := range subscribersList {
		processorList := make([]MessageProcessor, 0)

		// Add base logging processors
		//processorList = append(processorList, WithMessageLogging(subscriber.ChatID, msg))

		// Conditionally append db processors
		if ayah, ok := msg.(*models.Ayah); ok {
			processorList = append(processorList, WithDBPersistence(rs, &models.OutgoingMessage{
				ReceiverChatID: subscriber.ChatID,
				AyahID:         &ayah.ID,
			}))
		}

		err = rs.Bot.SendMessage(rs, msg, subscriber.ChatID, processorList...)
		if err != nil {
			log.Println(fmt.Sprintf("error while sending msg to : %v, chatID : %v", subscriber.UserName, subscriber.ChatID), err)
		}
	}
	return nil
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
