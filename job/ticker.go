package job

import (
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/controller"
	"time"
)

func StartTicker(rs *controller.Resource) {
	ticker := time.NewTicker(viper.GetDuration("ticker.duration"))

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go func() {
				log.Println("######################### Started Fetching #########################")
				ayah := rs.FetchNewVerse()
				err := rs.PublishToSubscribers(ayah)
				if err != nil {
					log.Println("err while publishing : ", err)
				}
				log.Println("######################### Finished Fetching #########################\n\n")
			}()
		}
	}
}
