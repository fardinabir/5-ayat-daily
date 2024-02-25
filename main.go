package main

import (
	"one-minute-quran/controller"
	"one-minute-quran/job"
)

func main() {
	rs := controller.NewResource()
	go job.StartTicker(rs)
	rs.ServeBot()
}
