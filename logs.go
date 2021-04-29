package main

import (
	"fmt"
	"sync"
	"time"
)

var logsQueue = []string{}
var logsMutex sync.Mutex

func getLogs() []string {
	return logsQueue
}

func newLog(v ...interface{}) {
	logStr := time.Now().Local().Format("Jan 2 15:04:05 - ") + fmt.Sprint(v...)
	fmt.Println(logStr)
	logsMutex.Lock()
	logsQueue = append(logsQueue, logStr)
	logsMutex.Unlock()
}
