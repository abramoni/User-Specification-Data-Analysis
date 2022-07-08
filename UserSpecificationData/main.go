package main

import (
	"fmt"
	"presentation/mypackages"
	//"time"
)

const (
	timezoneJakarta = "Asia"
)

func main() {

	var count int64
	count = 100

	var start_time, end_time, timestamp int64
	start_time = 1656453600
	end_time = 1656489600

	mypackages.UserRequestForData(start_time, end_time)

	fmt.Println("Alerts occured during this period are : ")

	for timestamp = start_time; timestamp <= end_time; timestamp += 3600 {

		mypackages.RetriveFronElasticDb(timestamp)
		mypackages.Examine(timestamp, count)

	}

}
