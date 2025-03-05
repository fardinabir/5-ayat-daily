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
		log.Println("error sending telegram message:", err)
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
			userName := update.Message.From.FirstName + " " + update.Message.From.LastName
			log.Println("Message from user : ", userName, " chatID : ", chatID, " Message Text : ", messageText)

			command := update.Message.Command()
			args := update.Message.CommandArguments()

			// process command
			switch command {
			case "start":
				t.handleStart(rs, chatID, userName)
			case "subscribe":
				t.handleSubscribe(rs, chatID, userName)
			case "unsubscribe":
				t.handleUnsubscribe(rs, chatID, userName)
			case "next":
				t.fetchNextVerse(rs, chatID)
			case "previous":
				t.fetchPreviousVerse(rs, chatID)
			case "random":
				t.fetchRandomVerse(rs, chatID)
			case "get_ayat":
				t.GetAyah(rs, args, chatID)
			case "insert_preferred":
				t.SavePreference(rs, args, chatID)
			case "broadcast":
				t.BroadCastMessage(rs, args, chatID)
			default:
				t.handleInvalidCommand(rs, chatID)
				command = ""
			}

			rs.Store.SaveIncomingMessage(&models.IncomingMessage{
				ChatID:         chatID,
				UserName:       userName,
				MessageText:    messageText,
				MessageCommand: command,
			})
		}
	}
}

func (t *tgBot) handleStart(rs *Resource, chatID, userName string) error {
	// Send a welcome message
	err := t.SendMessage(rs, fmt.Sprintf("Hello %v, Welcome!\n\nTo subscribe the channel click or type:\n/subscribe\n", userName), chatID, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (t *tgBot) SavePreference(rs *Resource, args, chatID string) error {
	argFields := strings.Fields(args)
	adminID := viper.GetString("telegram.adminID")
	if chatID != adminID {
		return fmt.Errorf("unauthorized request to a command")
	}

	if len(argFields) != 2 {
		if err := t.SendMessage(rs, "Please follow this format:\n/get_ayat <suraNo> <ayatNo>", chatID, nil); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong command format")
	}
	sId, _ := strconv.Atoi(argFields[0])
	vId, _ := strconv.Atoi(argFields[1])

	ayah, err := rs.Store.GetAyahSuraVerse(sId, vId)
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

func (t *tgBot) BroadCastMessage(rs *Resource, args, chatID string) error {
	adminID := viper.GetString("telegram.adminID")
	if chatID != adminID {
		return fmt.Errorf("unauthorized request to a command")
	}
	return rs.PublishToSubscribers(nil, args)
}

func (t *tgBot) GetAyah(rs *Resource, args string, chatID string) error {
	argFields := strings.Fields(args)
	if len(argFields) != 2 {
		if err := t.SendMessage(rs, "Please follow this format:\n/get_ayat <suraNo> <ayatNo>", chatID, nil); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong command format")
	}
	sId, _ := strconv.Atoi(argFields[0])
	vId, _ := strconv.Atoi(argFields[1])
	ayah, err := rs.Store.GetAyahSuraVerse(sId, vId)
	if err != nil {
		if err := t.SendMessage(rs, "Couldn't fetch requested ayat, Please follow this format:\n/get_ayat <suraNo> <ayatNo>", chatID, nil); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong sura-ayat format")
	}
	ayahText := FormatAyahText(ayah)

	if err := t.SendMessage(rs, ayahText, chatID, &ayah.ID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleSubscribe(rs *Resource, chatID, userName string) error {
	sub, err := rs.Store.GetSubscriber(chatID)
	if sub != nil {
		if err := t.SendMessage(rs, "You are already subscribed!", chatID, nil); err != nil {
			return fmt.Errorf("failed to send subscribe message: %w", err)
		}
		return nil
	}

	subscriber := &models.Subscriber{
		ChatID:   chatID,
		UserName: userName,
		Status:   "active",
		Channel:  "telegram",
	}
	err = rs.Store.Create(subscriber)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			// means user is in unsubscribed status, now tries to subscribe
			// subscriber.ChatID = subscriber.ChatID + "_old"
			err = rs.Store.HardDeleteSubscriber(chatID)
			log.Println("Found old subscription deleting current one, trying again with : ", subscriber.ChatID)
			t.handleSubscribe(rs, chatID, userName)
		}
		return nil
	}
	if err := t.SendMessage(rs, "Your subscription is greatly appreciated – thanks for joining us!\nHave your first Ayat of this day...", chatID, nil); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	t.fetchRandomVerse(rs, chatID)

	err = t.SendMessage(rs, fmt.Sprintf("Available Commands :\n\n"+
		"/subscribe - to get subscribed and receive updates daily\n"+
		"/next - to get the next ayah\n"+
		"/previous - to get the previous ayah\n"+
		"/random - to get a random ayah\n"+
		"/get_ayat <suraNo> <ayatNo> - to get a specific ayat from given sura\n"+
		"/unsubscribe - to unsubscribe daily updates"), chatID, nil)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	if err := t.SendMessage(rs, fmt.Sprintf("%v joined the channel", userName), adminID, nil); err != nil {
		return fmt.Errorf("failed to send admin notification: %w", err)
	}
	return nil
}

func (t *tgBot) handleUnsubscribe(rs *Resource, chatID, userName string) error {
	err := rs.Store.DeleteSubscriber(chatID)

	if err != nil {
		return err
	}
	if err := t.SendMessage(rs, "You are unsubscribed successfully, thanks for being part of the community.", chatID, nil); err != nil {
		return fmt.Errorf("failed to send unsubscribe message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	if err := t.SendMessage(rs, fmt.Sprintf("%v unsubscribed from the channel", userName), adminID, nil); err != nil {
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
