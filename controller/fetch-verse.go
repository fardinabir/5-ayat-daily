package controller

import (
	"five-ayat-daily/models"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"time"
)

func (rs *Resource) FetchNewVerse() *models.Ayah {
	vc := viper.GetString("operation.verse_preference")
	if vc == "enable" {
		pv, err := rs.Store.GetPreferredVerse()
		if err == nil {
			fmt.Println("----------------verse preference enabled---------------")
			ayah, _ := rs.Store.GetAyah(pv.VerseId)
			pv.Status = "sent"
			rs.Store.SavePreferredVerse(pv)
			return ayah
		}
	}
	return rs.FetchRandomVerse()
}

func (rs *Resource) FetchRandomVerse() *models.Ayah {
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

func (rs *Resource) FetchPreviousVerse(ayahId int) *models.Ayah {
	nextAyah := ayahId - 1
	if nextAyah == 0 {
		nextAyah = 6236
	}
	log.Println("-------- fetching previous verse : ", nextAyah)
	ayah, _ := rs.Store.GetAyah(nextAyah)
	return ayah
}
