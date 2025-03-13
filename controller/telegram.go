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

func (t *tgBot) SendMessage(rs *Resource, msg Message, chatID string, processors ...MessageProcessor) error {
	chatId, _ := strconv.Atoi(chatID)
	msgCfg := tgbotapi.NewMessage(int64(chatId), msg.GetContent())
	_, err := t.API.Send(msgCfg)
	if err != nil {
		log.Println("error sending telegram message:", err)
		return err
	}

	// Execute all message processors
	for _, process := range processors {
		if err := process(); err != nil {
			log.Printf("processor error: %v", err)
			return err
		}
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
	gMsg := models.GeneralMessage{Message: fmt.Sprintf("Hello %v, Welcome!\n\nTo subscribe the channel click or type:\n/subscribe\n", userName)}
	err := t.SendMessage(rs, gMsg, chatID)
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
		gMsg := models.GeneralMessage{Message: "Please follow this format:\n/get_ayat <suraNo> <ayatNo>"}
		if err := t.SendMessage(rs, gMsg, chatID); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong command format")
	}
	sId, _ := strconv.Atoi(argFields[0])
	vId, _ := strconv.Atoi(argFields[1])

	ayah, err := rs.Store.GetAyahSuraVerse(sId, vId)
	if err != nil {
		gMsg := models.GeneralMessage{Message: "Failed to get ayah with this combination"}
		if err := t.SendMessage(rs, gMsg, chatID); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}
	rs.Store.SavePreferredVerse(&models.VersePreference{
		VerseId: int(ayah.ID),
	})
	ayahText := models.GeneralMessage{Message: ayah.GetContent() + "\n\n --------- Saved as preferred verse -------- "}

	if err := t.SendMessage(rs, ayahText, chatID); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) BroadCastMessage(rs *Resource, args, chatID string) error {
	adminID := viper.GetString("telegram.adminID")
	if chatID != adminID {
		return fmt.Errorf("unauthorized request to a command")
	}

	return rs.PublishToSubscribers(models.GeneralMessage{Message: args})
}

func (t *tgBot) GetAyah(rs *Resource, args string, chatID string) error {
	argFields := strings.Fields(args)
	if len(argFields) != 2 {
		gMsg := models.GeneralMessage{Message: "Please follow this format:\n/get_ayat <suraNo> <ayatNo>"}
		if err := t.SendMessage(rs, gMsg, chatID); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong command format")
	}
	sId, _ := strconv.Atoi(argFields[0])
	vId, _ := strconv.Atoi(argFields[1])
	ayah, err := rs.Store.GetAyahSuraVerse(sId, vId)
	if err != nil {
		gMsg := models.GeneralMessage{Message: "Couldn't fetch requested ayat, Please follow this format:\n/get_ayat <suraNo> <ayatNo>"}
		if err := t.SendMessage(rs, gMsg, chatID); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		return fmt.Errorf("wrong sura-ayat format")
	}
	outgoingAyah := models.OutgoingMessage{
		ReceiverChatID: chatID,
		AyahID:         &ayah.ID,
	}

	if err := t.SendMessage(rs, ayah, chatID, WithDBPersistence(rs, &outgoingAyah)); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleSubscribe(rs *Resource, chatID, userName string) error {
	sub, err := rs.Store.GetSubscriber(chatID)
	if sub != nil {
		gMsg := models.GeneralMessage{Message: "You are already subscribed!"}
		if err := t.SendMessage(rs, gMsg, chatID); err != nil {
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
	gMsg := models.GeneralMessage{Message: "Your subscription is greatly appreciated – thanks for joining us!\nHave your first Ayat of this day..."}
	if err := t.SendMessage(rs, gMsg, chatID); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	t.fetchRandomVerse(rs, chatID)

	gMsg = models.GeneralMessage{Message: fmt.Sprintf("Available Commands :\n\n" +
		"/subscribe - to get subscribed and receive updates daily\n" +
		"/next - to get the next ayah\n" +
		"/previous - to get the previous ayah\n" +
		"/random - to get a random ayah\n" +
		"/get_ayat <suraNo> <ayatNo> - to get a specific ayat from given sura\n" +
		"/unsubscribe - to unsubscribe daily updates")}
	err = t.SendMessage(rs, gMsg, chatID)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	gMsg = models.GeneralMessage{Message: fmt.Sprintf("%v joined the channel", userName)}
	if err := t.SendMessage(rs, gMsg, adminID); err != nil {
		return fmt.Errorf("failed to send admin notification: %w", err)
	}
	return nil
}

func (t *tgBot) handleUnsubscribe(rs *Resource, chatID, userName string) error {
	err := rs.Store.DeleteSubscriber(chatID)

	if err != nil {
		return err
	}
	gMsg := models.GeneralMessage{Message: "You are unsubscribed successfully, thanks for being part of the community."}
	if err := t.SendMessage(rs, gMsg, chatID); err != nil {
		return fmt.Errorf("failed to send unsubscribe message: %w", err)
	}

	adminID := viper.GetString("telegram.adminID")
	gMsg = models.GeneralMessage{Message: fmt.Sprintf("%v unsubscribed from the channel", userName)}
	if err := t.SendMessage(rs, gMsg, adminID); err != nil {
		return fmt.Errorf("failed to send admin notification: %w", err)
	}
	return nil
}

func (t *tgBot) fetchNextVerse(rs *Resource, chatID string) error {
	lastMessage, err := rs.Store.GetLastOutgoingAyah(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gMsg := models.GeneralMessage{Message: "Subscribe to get messages, click or type /subscribe"}
			if err := t.SendMessage(rs, gMsg, chatID); err != nil {
				return fmt.Errorf("failed to send subscribe message: %w", err)
			}
		}
		log.Println("error while getting last outgoing message : ", err)
		return err
	}
	ayah := rs.FetchNextVerse(int(*lastMessage.AyahID))
	outgoingAyah := models.OutgoingMessage{
		ReceiverChatID: chatID,
		AyahID:         &ayah.ID,
	}

	if err := t.SendMessage(rs, ayah, chatID, WithDBPersistence(rs, &outgoingAyah)); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) fetchPreviousVerse(rs *Resource, chatID string) error {
	lastMessage, err := rs.Store.GetLastOutgoingAyah(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gMsg := models.GeneralMessage{Message: "Subscribe to get messages, click or type /subscribe"}
			if err := t.SendMessage(rs, gMsg, chatID); err != nil {
				return fmt.Errorf("failed to send subscribe message: %w", err)
			}
		}
		log.Println("error while getting last outgoing message : ", err)
		return err
	}
	ayah := rs.FetchPreviousVerse(int(*lastMessage.AyahID))
	outgoingAyah := models.OutgoingMessage{
		ReceiverChatID: chatID,
		AyahID:         &ayah.ID,
	}

	if err := t.SendMessage(rs, ayah, chatID, WithDBPersistence(rs, &outgoingAyah)); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) fetchRandomVerse(rs *Resource, chatID string) error {
	ayah := rs.FetchRandomVerse()
	outgoingAyah := models.OutgoingMessage{
		ReceiverChatID: chatID,
		AyahID:         &ayah.ID,
	}

	if err := t.SendMessage(rs, ayah, chatID, WithDBPersistence(rs, &outgoingAyah)); err != nil {
		return fmt.Errorf("failed to send ayah message: %w", err)
	}
	return nil
}

func (t *tgBot) handleInvalidCommand(rs *Resource, chatID string) error {
	gMsg := models.GeneralMessage{Message: "Invalid command or format, type '/' to see the available commands"}
	return t.SendMessage(rs, gMsg, chatID, nil)
}
