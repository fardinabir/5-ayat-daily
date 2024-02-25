package models

import (
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

func (t *tgBot) ServeBot() interface{} {
	// Update Config From TgBOT
	updateConfig := tgbotapi.NewUpdate(0)
	updates := t.API.GetUpdatesChan(updateConfig)
	return updates
}
