package job

import (
	"log"
	"one-minute-quran/controller"
	"one-minute-quran/controller/verse-loader"
	"time"
)

func StartTicker(rs *controller.Resource) {
	ticker := time.NewTicker(20 * time.Second)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("######################### Started Fetching #########################")
			verse := verse_loader.FetchNewVerse()
			rs.PublishToSubscribers(verse)
			log.Println("######################### Finished Fetching #########################\n\n")
		}
	}
}
