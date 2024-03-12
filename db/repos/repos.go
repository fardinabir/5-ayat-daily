package repos

import (
	"gorm.io/gorm"
	"log"
	"one-minute-quran/db"
	"one-minute-quran/models"
)

type Store struct {
	DB *gorm.DB
}

var storeInstance *Store

func NewSubsStore() *Store {
	if storeInstance != nil {
		return storeInstance
	}
	storeInstance = &Store{DB: db.ConnectDB()}
	return storeInstance
}

func (s *Store) Save(sd *models.Subscriber) error {
	res := s.DB.Save(sd)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}

func (s *Store) GetSubscriber(chatID string) (*models.Subscriber, error) {
	var sd models.Subscriber
	res := s.DB.Model(&models.Subscriber{}).Where("chat_id = ?", chatID).First(&sd)
	if res.Error != nil {
		log.Println("Error while getting subscriber in db", res.Error)
		return nil, res.Error
	}
	return &sd, nil
}

func (s *Store) GetAllSubscribers() ([]models.Subscriber, error) {
	var sd []models.Subscriber
	res := s.DB.Model(&models.Subscriber{}).Where("status = ? ", "active").Find(&sd)
	if res.Error != nil {
		log.Println("Error while getting subscriber in db", res.Error)
		return nil, res.Error
	}
	return sd, nil
}

func (s *Store) GetAyah(id int) (*models.Ayah, error) {
	var ay models.Ayah
	res := s.DB.First(&ay, id)
	if res.Error != nil {
		log.Println("Error while getting ayah in db", res.Error)
		return nil, res.Error
	}
	return &ay, nil
}

func (s *Store) SaveOutgoingMessage(ogMsg *models.OutgoingMessage) error {
	res := s.DB.Save(ogMsg)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}

func (s *Store) SaveIncomingMessage(inMsg *models.IncomingMessage) error {
	res := s.DB.Save(inMsg)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}

func (s *Store) GetLastOutgoingAyah(receiverChatID string) (*models.OutgoingMessage, error) {
	var om models.OutgoingMessage
	res := s.DB.Model(&models.OutgoingMessage{}).Where("receiver_chat_id = ? AND ayah_id IS NOT NULL", receiverChatID).Order("id desc").Limit(1).First(&om)
	if res.Error != nil {
		log.Println("Error while getting ayah in db", res.Error)
		return nil, res.Error
	}
	return &om, nil
}
