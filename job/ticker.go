package job

import (
	"log"
	"one-minute-quran/controller"
	verse_loader "one-minute-quran/controller/verse-loader"
	"time"
)

func StartTicker(rs *controller.Resource) {
	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("######################### Started Fetching #########################")
				verse := verse_loader.FetchNewVerse()
				err := rs.PublishToSubscribers(verse)
				if err != nil {
					log.Println("err while publishing : ", err)
				}
				log.Println("######################### Finished Fetching #########################\n\n")
			}()
		}
	}
}
