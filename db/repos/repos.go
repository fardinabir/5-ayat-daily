package repos

import (
	"gorm.io/gorm"
	"log"
	"one-minute-quran/db"
	"one-minute-quran/models"
)

type SubscriberStore struct {
	DB *gorm.DB
}

var subsStore *SubscriberStore

func NewSubsStore() *SubscriberStore {
	if subsStore != nil {
		return subsStore
	}
	subsStore = &SubscriberStore{DB: db.ConnectDB()}
	return subsStore
}

func (s *SubscriberStore) Save(sd *models.Subscriber) error {
	res := s.DB.Save(sd)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}

func (s *SubscriberStore) GetSubscriber(chatID string) (*models.Subscriber, error) {
	var sd models.Subscriber
	res := s.DB.Model(&models.Subscriber{}).Where("chat_id = ?", chatID).First(&sd)
	if res.Error != nil {
		log.Println("Error while getting subscriber in db", res.Error)
		return nil, res.Error
	}
	return &sd, nil
}

func (s *SubscriberStore) GetAllSubscribers() ([]models.Subscriber, error) {
	var sd []models.Subscriber
	res := s.DB.Model(&models.Subscriber{}).Where("status = ? ", "active").Find(&sd)
	if res.Error != nil {
		log.Println("Error while getting subscriber in db", res.Error)
		return nil, res.Error
	}
	return sd, nil
}

func (s *SubscriberStore) GetAyah(id int) (*models.Ayah, error) {
	var ay models.Ayah
	res := s.DB.First(&ay, id)
	if res.Error != nil {
		log.Println("Error while getting ayah in db", res.Error)
		return nil, res.Error
	}
	return &ay, nil
}

func (s *SubscriberStore) SaveOutgoingMessage(ogMsg *models.OutgoingMessage) error {
	res := s.DB.Save(ogMsg)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}

func (s *SubscriberStore) SaveIncomingMessage(inMsg *models.IncomingMessage) error {
	res := s.DB.Save(inMsg)
	if res.Error != nil {
		log.Println("Error while creating entry in db", res.Error)
		return res.Error
	}
	return nil
}
