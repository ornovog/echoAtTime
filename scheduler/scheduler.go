package scheduler

import (
	"echoAtTime/storageHandler"
	"fmt"
	"time"
)

type scheduler struct{
	storageHandler storageHandler.StorageInterface
}

func NewScheduler() scheduler {
	return scheduler{}
}

func (s *scheduler) Init(storageHandler storageHandler.StorageInterface){
	s.storageHandler = storageHandler
	go s.readAndScheduleMessages()
}

func(s scheduler) readAndScheduleMessages(){
	for{
		message := s.storageHandler.GetNextMessage()
		go s.handleMessage(message)
	}
}

func (s scheduler) handleMessage(m storageHandler.Message){
	text := m.Text
	delayInSeconds := extractDelayTimeInSeconds(m.Unix)

	printAfterDelay(delayInSeconds, text)
}

func extractDelayTimeInSeconds(unix int64)int64 {
	delayInSeconds := unix - nowUnix()
	if delayInSeconds < 0 {
		delayInSeconds = 0
	}
	return delayInSeconds
}

func printAfterDelay(delayInSecond int64, message string) {
	fmt.Println(delayInSecond)
	<-time.After(time.Duration(delayInSecond) * time.Second)
	fmt.Println(message)
}

func nowUnix()int64{
	return time.Now().Add(2*time.Hour).Unix()//I got a delay of 2 hours in time.Now()
}



