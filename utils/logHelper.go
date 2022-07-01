package utils

import (
	"fmt"
	"time"
)

func SleepLog(msg string, totalDuration time.Duration) {
	interval := time.Second * 1
	currentTime := time.Now()
	endTime := currentTime.Add(totalDuration)
	for currentTime.Before(endTime) {
		diff := endTime.Sub(currentTime)
		second := int(diff.Seconds())
		fmt.Printf("msg : %s, %d secs remaining...\n", msg, second)
		time.Sleep(interval)
		currentTime = currentTime.Add(interval)
	}
	fmt.Printf("msg : %s, done\n", msg)
}
