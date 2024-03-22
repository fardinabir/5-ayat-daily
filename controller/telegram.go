package controller

import (
	"five-ayat-daily/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

type tgBot struct {
	API *tgbotapi.BotAPI
}

var tgBotInstance *tgBot

func NewTgBot() (*tgBot, error) {
	if tgBotInstance != nil {
		return tgBotInstance, nil
	}
	token := viper.GetString("telegram.token")
	tgBotAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("error while bot initialization", err)
		return nil, err
	}
	tgBotInstance = &tgBot{tgBotAPI}
	return tgBotInstance, nil
}

func (t *tgBot) SendMessage(rs *Resource, message, chatID string, ayahId *uint) error {
	msg := &models.OutgoingMessage{
		ReceiverChatID: chatID,
		AyahID:         ayahId,
	}
	if ayahId == nil {
		msg.GeneralMessage = message
	}
	rs.Store.SaveOutgoingMessage(msg)

	chatId, _ := strconv.Atoi(chatID)
	msgCfg := tgbotapi.NewMessage(int64(chatId), message)
	_, err := t.API.Send(msgCfg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (t *tgBot) ServeBotAPI(rs *Resource) {
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
			command := strings.Split(messageText, " ")[0] // first word of messageText

			rs.Store.SaveIncomingMessage(&models.IncomingMessage{
				ChatID:         chatID,
				MessageText:    messageText,
				MessageCommand: command,
			})

			// process command
			if command == "/start" {
				// Send a welcome message
				t.handleStart(rs, chatID)
			} else if command == "/subscribe" {
				t.handleSubscribe(rs, chatID)
			} else if command == "/next" {
				t.fetchNextVerse(rs, chatID)
			} else if command == "/previous" {
				t.fetchPreviousVerse(rs, chatID)
			} else if command == "/random" {
				t.fetchRandomVerse(rs, chatID)
			} else {
				t.handleInvalidCommand(rs, chatID)
			}
		}
	}
}

func (t *tgBot) handleStart(rs *Resource, chatID string) error {
	// Send a welcome message
	return t.SendMessage(rs, "Hello User, Welcome!\n To subscribe the channel click or type:\n/subscribe", chatID, nil)
}

func (t *tgBot) handleSubscribe(rs *Resource, chatID string) error {
	rs.Store.Save(&models.Subscriber{
		ChatID:  chatID,
		Status:  "active",
		Channel: "telegram",
	})

	if err := t.SendMessage(rs, "Now you are subscribed! Thanks!", chatID, nil); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	if err := t.SendMessage(rs, fmt.Sprintf("%v joined the channel", chatID), adminID, nil); err != nil {
		return fmt.Errorf("failed to send admin notification: %w", err)
	}
	return nil
}

func (t *tgBot) fetchNextVerse(rs *Resource, chatID string) error {
	lastMessage, err := rs.Store.GetLastOutgoingAyah(chatID)
	if err != nil {
		log.Println("error while getting last outgoing message : ", err)
		return err
	}
	ayah := rs.FetchNextVerse(int(*lastMessage.AyahID))
	ayahText := FormatAyahText(ayah)

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) fetchPreviousVerse(rs *Resource, chatID string) error {
	lastMessage, err := rs.Store.GetLastOutgoingAyah(chatID)
	if err != nil {
		log.Println("error while getting last outgoing message : ", err)
		return err
	}
	ayah := rs.FetchPreviousVerse(int(*lastMessage.AyahID))
	ayahText := FormatAyahText(ayah)

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) fetchRandomVerse(rs *Resource, chatID string) error {
	ayah := rs.FetchNewVerse()
	ayahText := FormatAyahText(ayah)

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleInvalidCommand(rs *Resource, chatID string) error {
	return t.SendMessage(rs, "Invalid command or format, type '/' to see the available commands", chatID, nil)
}
