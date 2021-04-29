package main

import (
	"sync"
	"time"
)

// globalState is key, value pair used to store temp state and can be accessed anytime by any func
var globalState = map[string](string){}
var stateMutex sync.Mutex
var lastUpdated = time.Now()

func getState(key string) string {
	return globalState[key]
}

func setState(key string, value string) {
	stateMutex.Lock()
	globalState[key] = value
	lastUpdated = time.Now()
	stateMutex.Unlock()
}
