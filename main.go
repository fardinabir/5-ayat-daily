package main

import (
	"five-ayat-daily/controller"
	"five-ayat-daily/job"
	"log"
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
