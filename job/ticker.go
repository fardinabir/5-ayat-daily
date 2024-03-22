package job

import (
	"five-ayat-daily/controller"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

func StartTicker(rs *controller.Resource, wg *sync.WaitGroup) {
	ticker := time.NewTicker(viper.GetDuration("ticker.duration"))

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			trigger := compareTimes()
			if trigger {
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
}

func compareTimes() bool {
	dhakaTime, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return false
	}

	// Use the loaded time zone
	currentTime := time.Now().In(dhakaTime).Format("15:04")
	currentTimeFlat, _ := time.Parse("15:04", currentTime)
	// Read trigger time values from the configuration file
	triggerTimes := viper.GetStringSlice("trigger_times")

	// Check matching time
	for _, triggerTime := range triggerTimes {
		parsedTime, err := time.Parse("15:04", triggerTime)
		if err != nil {
			log.Println("failed to parse time : ", err)
			return false
		}

		// check close match for given time and cur time
		diff := parsedTime.Sub(currentTimeFlat)
		if diff < 5*time.Minute && diff >= 0 {
			log.Println("Caught nearest trigger time, diff is : ", diff)
			return true
		}
		//log.Println("currentTime : ", currentTime, " parsedTime : ", parsedTime, " diff : ", diff)
	}
	return false
}
