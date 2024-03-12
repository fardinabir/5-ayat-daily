package controller

import (
	"math/rand"
	"one-minute-quran/models"
	"time"
)

func (rs *Resource) FetchNewVerse() *models.Ayah {
	rand.Seed(time.Now().UnixNano())
	// Generate a random number between 6236
	randomNumber := rand.Intn(6236) + 1 // Adding 1 to include 6236 in the range

	ayah, _ := rs.SubsStore.GetAyah(randomNumber)
	return ayah
}
