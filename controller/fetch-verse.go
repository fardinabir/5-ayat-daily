package controller

import (
	"log"
	"math/rand"
	"one-minute-quran/models"
	"time"
)

func (rs *Resource) FetchNewVerse() *models.Ayah {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(6236) + 1 // Adding 1 to include 6236 in the range

	ayah, _ := rs.Store.GetAyah(randomNumber)
	return ayah
}

func (rs *Resource) FetchNextVerse(ayahId int) *models.Ayah {
	nextAyah := (ayahId + 1) % 6236
	if nextAyah == 0 {
		nextAyah = 6236
	}
	log.Println("-------- fetching next verse : ", nextAyah)
	ayah, _ := rs.Store.GetAyah(nextAyah)
	return ayah
}
