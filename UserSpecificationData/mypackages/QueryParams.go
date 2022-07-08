package mypackages

import (
	"fmt"
	"strconv"
	"time"
)

func PresentTimeStamp() (timestamp_now int64) {
	tNow := time.Now()
	timestamp_now = tNow.Unix()
	return
}

func TimeStampGneerator(timeperiod string, unit string) (startTime string, endTime string) {

	var unix_start_time, unix_end_time int64

	if timeperiod == "day" {
		unix_end_time = Present_time_stamp
		unix_start_time = unix_end_time - 3600
	}

	if unit == "M" {
		unix_start_time = unix_start_time * 1000
		unix_end_time = unix_end_time * 1000
	}
	if unit == "U" {
		unix_start_time = unix_start_time * 1000000
		unix_end_time = unix_end_time * 1000000
	}
	s := strconv.FormatInt(unix_start_time, 10)
	e := strconv.FormatInt(unix_end_time, 10)
	return s, e
}

func PrintTime(timestamp int64) {
	timeT := time.Unix(timestamp, 0)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := timeT.In(loc)
	time_start := timestamp - 3600
	start_timeT := time.Unix(time_start, 0)
	start_now := start_timeT.In(loc)

	fmt.Println("During the period from ", start_now, " to ", now)
}
