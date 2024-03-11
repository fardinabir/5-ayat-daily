package models

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type tgBot struct {
	API *tgbotapi.BotAPI
}

//var tgBotAPI *tgbotapi.BotAPI

func NewTgBot() *tgBot {
	token := viper.GetString("telegram.token")
	tgBotAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil
	}
	return &tgBot{tgBotAPI}
}

func (t *tgBot) SendMessage(message, chatID string) error {
	chatId, _ := strconv.Atoi(chatID)
	msgCfg := tgbotapi.NewMessage(int64(chatId), message)
	_, err := t.API.Send(msgCfg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (t *tgBot) ServeBot(subsChan chan<- Subscriber) {
	// Update Config From TgBOT
	updateConfig := tgbotapi.NewUpdate(0)
	updates := t.API.GetUpdatesChan(updateConfig)

	// Process messages from updates chan
	for update := range updates {
		if update.Message != nil {
			// Get the message text and chat ID
			messageText := update.Message.Text
			chatID := strconv.Itoa(int(update.Message.Chat.ID))
			log.Println("Message from chatID : ", chatID, "Message Text : ", messageText)

			if messageText == "/start" {
				// Send a welcome message
				t.handleStart(chatID)
			} else if messageText == "/subscribe" {
				t.handleSubscribe(chatID, subsChan)
			}
		}
	}
	close(subsChan)
}

func (t *tgBot) handleStart(chatID string) error {
	// Send a welcome message
	return t.SendMessage("Hello User, Welcome!\n To subscribe the channel click or type:\n/subscribe", chatID)
}

func (t *tgBot) handleSubscribe(chatID string, subsChan chan<- Subscriber) error {
	sub := Subscriber{
		ChatID:  chatID,
		Status:  "active",
		Channel: "telegram",
	}
	// sends to newSubscriber Channel
	subsChan <- sub

	if err := t.SendMessage("Now you are subscribed to one-minute-quran! Thanks!", chatID); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	if err := t.SendMessage(fmt.Sprintf("%v joined the channel", chatID), adminID); err != nil {
		return fmt.Errorf("failed to send admin notification: %w", err)
	}

	return nil
}
