package db

import "one-minute-quran/models"

func getModels() []interface{} {
	var Models []interface{}
	Models = append(Models, &models.Subscriber{})
	Models = append(Models, &models.IncomingMessage{})
	Models = append(Models, &models.OutgoingMessage{})
	Models = append(Models, &models.Ayah{})
	Models = append(Models, &models.Category{})

	return Models
}
