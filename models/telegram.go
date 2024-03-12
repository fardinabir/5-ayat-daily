package models

//
//import (
//	"fmt"
//	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/spf13/viper"
//	"log"
//	"one-minute-quran/controller"
//	"strconv"
//	"strings"
//)
//
//type tgBot struct {
//	API *tgbotapi.BotAPI
//}
//
////var tgBotAPI *tgbotapi.BotAPI
//
//var tgBotInstance *tgBot
//
//func NewTgBot() *tgBot {
//	if tgBotInstance != nil {
//		return tgBotInstance
//	}
//	token := viper.GetString("telegram.token")
//	tgBotAPI, err := tgbotapi.NewBotAPI(token)
//	if err != nil {
//		log.Println("error while bot initialization", err)
//		return nil
//	}
//	tgBotInstance = &tgBot{tgBotAPI}
//	return tgBotInstance
//}
//
//func (t *tgBot) SendMessage(message, chatID string) error {
//	chatId, _ := strconv.Atoi(chatID)
//	msgCfg := tgbotapi.NewMessage(int64(chatId), message)
//	_, err := t.API.Send(msgCfg)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	return nil
//}
//
//func (t *tgBot) ServeBot() {
//	// Update Config From TgBOT
//	updateConfig := tgbotapi.NewUpdate(0)
//	updates := t.API.GetUpdatesChan(updateConfig)
//
//	// Process messages from updates chan
//	for update := range updates {
//		if update.Message != nil {
//			// Get the message text and chat ID
//			messageText := update.Message.Text
//			chatID := strconv.Itoa(int(update.Message.Chat.ID))
//			log.Println("Message from chatID : ", chatID, "Message Text : ", messageText)
//			command := strings.Split(messageText, " ")[0] // first word of messageText
//
//			t.Rs.Store.SaveIncomingMessage(&IncomingMessage{
//				ChatID:         chatID,
//				MessageText:    messageText,
//				MessageCommand: command,
//			})
//
//			// process command
//			if command == "/start" {
//				// Send a welcome message
//				t.handleStart(chatID)
//			} else if command == "/subscribe" {
//				t.handleSubscribe(chatID)
//			} else if command == "/next" {
//				t.fetchNextVerse(chatID)
//			}
//		}
//	}
//}
//
//func (t *tgBot) handleStart(chatID string) error {
//	// Send a welcome message
//	return t.SendMessage("Hello User, Welcome!\n To subscribe the channel click or type:\n/subscribe", chatID)
//}
//
//func (t *tgBot) handleSubscribe(chatID string) error {
//	t.Rs.Store.Save(&Subscriber{
//		ChatID:  chatID,
//		Status:  "active",
//		Channel: "telegram",
//	})
//
//	if err := t.SendMessage("Now you are subscribed to one-minute-quran! Thanks!", chatID); err != nil {
//		return fmt.Errorf("failed to send subscribe message: %w", err)
//	}
//
//	adminID := viper.GetString("telegram.adminID")
//	if err := t.SendMessage(fmt.Sprintf("%v joined the channel", chatID), adminID); err != nil {
//		return fmt.Errorf("failed to send admin notification: %w", err)
//	}
//	return nil
//}
//
//func (t *tgBot) fetchNextVerse(chatID string) error {
//	ayah := t.Rs.FetchNewVerse()
//	ayahText := controller.FormatAyahText(ayah)
//
//	t.Rs.Store.SaveOutgoingMessage(&OutgoingMessage{
//		ReceiverType: RECEIVERTYPESINGLE,
//		AyahID:       ayah.ID,
//	})
//	if err := t.SendMessage(ayahText, chatID); err != nil {
//		return fmt.Errorf("failed to send ayah message: %w", err)
//	}
//	return nil
//}
