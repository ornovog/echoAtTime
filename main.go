package main

import (
	"echoAtTime/scheduler"
	"echoAtTime/storageHandler"
	"net/http"
	"time"
)

var messagesWriter chan storageHandler.Message

func main() {
	messagesWriter = make(chan storageHandler.Message)

	storageHandler := storageHandler.NewStorageHandler()
	storageHandler.Init(messagesWriter)

	scheduler := scheduler.NewScheduler()
	scheduler.Init(&storageHandler)

	http.HandleFunc("/echoAtTime", handleRequest)
	http.ListenAndServe(":8080", nil)
}


func handleRequest(_ http.ResponseWriter, request *http.Request){
	request.ParseForm()

	text := request.Form.Get("Text")
	timeStr := request.Form.Get("Time")

	m := extractMessage(text, timeStr)
	messagesWriter <- m
}

func extractMessage(text string, timeStr string) storageHandler.Message {
	var m storageHandler.Message
	m.Text = text
	m.Unix = timeStringToUnix(timeStr)

	return m
}

func timeStringToUnix(timeStr string)int64{
	layOut := "2006-01-02 15:04:05"
	timeStamp, _ := time.Parse(layOut, timeStr)
	return timeStamp.Unix()
}
