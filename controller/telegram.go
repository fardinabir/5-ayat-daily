package controller

import (
	"errors"
	"five-ayat-daily/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
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
			texts := strings.Split(messageText, " ") // first word of messageText
			command := texts[0]
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
			} else if command == "/insertPreferred" {
				t.SavePreference(rs, texts[1], texts[2], chatID)
			} else {
				t.handleInvalidCommand(rs, chatID)
				command = ""
			}

			rs.Store.SaveIncomingMessage(&models.IncomingMessage{
				ChatID:         chatID,
				MessageText:    messageText,
				MessageCommand: command,
			})
		}
	}
}

func (t *tgBot) handleStart(rs *Resource, chatID string) error {
	// Send a welcome message
	return t.SendMessage(rs, "Hello User, Welcome!\n To subscribe the channel click or type:\n/subscribe", chatID, nil)
}

func (t *tgBot) SavePreference(rs *Resource, surahId, verseId, chatID string) error {
	sId, _ := strconv.Atoi(surahId)
	vId, _ := strconv.Atoi(verseId)
	ayah, err := rs.Store.GetAyahSuraVerse(sId, vId)
	fmt.Println(ayah)
	if err != nil {
		if err := t.SendMessage(rs, "Failed to get ayah with this combination", chatID, nil); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}
	rs.Store.SavePreferredVerse(&models.VersePreference{
		VerseId: int(ayah.ID),
	})
	ayahText := FormatAyahText(ayah) + "\n\n --------- Saved as preferred verse -------- "

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleSubscribe(rs *Resource, chatID string) error {
	err := rs.Store.Create(&models.Subscriber{
		ChatID:  chatID,
		Status:  "active",
		Channel: "telegram",
	})

	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		if err := t.SendMessage(rs, "You are already subscribed!", chatID, nil); err != nil {
			return fmt.Errorf("failed to send subscribe message: %w", err)
		}
		return nil
	}
	t.fetchRandomVerse(rs, chatID)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := t.SendMessage(rs, "Subscribe to get messages, click or type /subscribe", chatID, nil); err != nil {
				return fmt.Errorf("failed to send subscribe message: %w", err)
			}
		}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := t.SendMessage(rs, "Subscribe to get messages, click or type /subscribe", chatID, nil); err != nil {
				return fmt.Errorf("failed to send subscribe message: %w", err)
			}
		}
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
	ayah := rs.FetchRandomVerse()
	ayahText := FormatAyahText(ayah)

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleInvalidCommand(rs *Resource, chatID string) error {
	return t.SendMessage(rs, "Invalid command or format, type '/' to see the available commands", chatID, nil)
}
