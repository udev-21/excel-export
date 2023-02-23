package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var currentSessions = make(map[string]Session)

var serverPort, fileServerPort, defaultSessionPurgeHour int
var debug bool

func main() {
	flag.IntVar(&serverPort, "port", 8000, "server port")
	flag.IntVar(&fileServerPort, "filePort", 8001, "file server port")
	flag.IntVar(&defaultSessionPurgeHour, "purgeHour", 24, "default session purge hour")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	log.Printf("staring...\nserverPort: %d\nfileServerPort: %d\n", serverPort, fileServerPort)
	go unusedSessionPurger()
	RunServer(fmt.Sprint(serverPort), fmt.Sprint(fileServerPort), debug)
}

func unusedSessionPurger() {
	hours, err := time.ParseDuration(fmt.Sprintf(
		"%vh",
		defaultSessionPurgeHour,
	))

	if err != nil {
		panic(err)
	}

	ticker := time.Tick(time.Hour * hours)
	for range ticker {
		for k, v := range currentSessions {
			if time.Since(v.LastTimeUsed) > time.Hour*24 {
				delete(currentSessions, k)
			}
		}
	}
}
