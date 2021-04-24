package main

import (
	"fmt"
	"log"
	"sync"
)

var logsQueue = []string{}
var logsMutex sync.Mutex

func GetLogs() []string {
	return logsQueue
}

func NewLog(v ...interface{}) {
	newlog := fmt.Sprint(v...)
	log.Println(newlog)
	logsMutex.Lock()
	logsQueue = append(logsQueue, newlog)
	logsMutex.Unlock()
}
