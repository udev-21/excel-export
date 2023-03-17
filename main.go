package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
)

var currentSessions = struct {
	s map[string]Session
	m *sync.Mutex
}{
	s: make(map[string]Session),
	m: &sync.Mutex{},
}

var serverPort, fileServerPort, defaultSessionPurgeHour int
var debug bool

func main() {
	flag.IntVar(&serverPort, "port", 8000, "server port")
	flag.IntVar(&fileServerPort, "filePort", 8001, "file server port")
	flag.IntVar(&defaultSessionPurgeHour, "purgeHour", 24, "default session purge hour")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	log.Printf("staring...\nserverPort: %d\nfileServerPort: %d\n", serverPort, fileServerPort)
	// go unusedSessionPurger()
	RunServer(fmt.Sprint(serverPort), fmt.Sprint(fileServerPort), debug)
}

// func unusedSessionPurger() {

// 	ticker := time.Tick(time.Hour * time.Duration(defaultSessionPurgeHour))
// 	for range ticker {
// 		for k, v := range currentSessions {
// 			if time.Since(v.LastTimeUsed) > time.Hour*time.Duration(defaultSessionPurgeHour) {
// 				delete(currentSessions, k)
// 			}
// 		}
// 	}
// }
