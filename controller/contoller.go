package controller

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/controller/interfaces"
	"one-minute-quran/db"
	"one-minute-quran/db/repos"
	"one-minute-quran/models"
	"strconv"
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
	updates := rs.Bot.ServeBot().(tgbotapi.UpdatesChannel)
	// Process messages from updates chan
	for update := range updates {
		if update.Message != nil {
			// Get the message text and chat ID
			messageText := update.Message.Text
			chatID := strconv.Itoa(int(update.Message.Chat.ID))
			log.Println("Message from chatID : ", chatID, "Message Text : ", messageText)

			if messageText == "/start" {
				// Send a welcome message
				err := rs.Bot.SendMessage("Hello User, Welcome!\n To subscribe the channel click or type:\n/subscribe", chatID)
				if err != nil {
					log.Println(err)
				}
			} else if messageText == "/subscribe" {
				sub := &models.Subscriber{
					ChatID:  chatID,
					Status:  "active",
					Channel: "telegram",
				}
				err := rs.SubsStore.Save(sub)
				if err != nil {
					log.Println("Cannot Save Subscriber Info: ", err)
					continue
				}
				err = rs.Bot.SendMessage("Now you are subscribed to one-minute-quran! Thanks!", chatID)
				if err != nil {
					log.Println("failed to send message : ", err)
				}

				adminId := viper.GetString("telegram.adminID")
				err = rs.Bot.SendMessage(fmt.Sprintf("%v joined the channel", chatID), adminId)
				if err != nil {
					log.Println("Admin approve request msg: ", err)
				}
			}
		}
	}
}

func LoadFromConfig() {
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}
