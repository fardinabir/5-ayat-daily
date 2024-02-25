package repos

import (
	"gorm.io/gorm"
	"log"
	"one-minute-quran/models"
)

type SubscriberStore struct {
	DB *gorm.DB
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
