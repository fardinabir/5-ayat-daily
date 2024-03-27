package db

import "five-ayat-daily/models"

func getModels() []interface{} {
	var Models []interface{}
	Models = append(Models, &models.Subscriber{})
	Models = append(Models, &models.IncomingMessage{})
	Models = append(Models, &models.OutgoingMessage{})
	Models = append(Models, &models.Ayah{})
	Models = append(Models, &models.Category{})
	Models = append(Models, &models.VersePreference{})

	return Models
}
