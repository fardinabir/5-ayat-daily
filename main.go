package main

import (
	"log"
	"one-minute-quran/controller"
	"one-minute-quran/job"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	rs := controller.NewResource()
	if rs == nil {
		log.Println("error while resource initialization, exiting")
		return
	}
	wg.Add(1)
	go job.StartTicker(rs, &wg)
	rs.ServeBot()
	wg.Wait()
}
