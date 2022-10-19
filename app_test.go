package main

import (
	"log"
	"testing"
	"time"
)

var count int

func TestStopwatchManager(t *testing.T) {
	swMng := NewStopwatchManager(stopwatchCallback)
	defer swMng.StopAll()

	swMng.Start(1)
	time.Sleep(2 * time.Second)
	if count > 3 || count < 1 {
		t.Error("timer not at 10 after 10 seconds: ", count)
	}
	t.Log("stop timer 1")
	swMng.Stop(1)
	t.Log("wait stop all")
}

func stopwatchCallback(id int32, time time.Duration) {
	if id == 1 {
		log.Println("clbk called: ", id, time)
		count++
	}
}
