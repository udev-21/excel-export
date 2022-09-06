package main

import (
	"flag"
	"fmt"
	"log"
)

var currentSessions = make(map[string]Session)

var serverPort, fileServerPort int
var debug bool

func main() {
	flag.IntVar(&serverPort, "port", 8000, "server port")
	flag.IntVar(&fileServerPort, "filePort", 8001, "file server port")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	log.Printf("staring...\nserverPort: %d\nfileServerPort: %d\n", serverPort, fileServerPort)
	RunServer(fmt.Sprint(serverPort), fmt.Sprint(fileServerPort), debug)
}
