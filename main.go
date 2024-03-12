package main

import (
	"one-minute-quran/controller"
	"one-minute-quran/job"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	rs := controller.NewResource()
	if rs == nil {
		return
	}
	wg.Add(1)
	go job.StartTicker(rs, &wg)
	rs.ServeBot()
	wg.Wait()
}
