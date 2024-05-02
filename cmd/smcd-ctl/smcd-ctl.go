package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Default().Fatalln("Please provide a sub-command: start,stop,status")
	}

	switch os.Args[1] {
	case "start":

	case "stop":

	default:
		log.Default().Fatalf("Invalid sub-command: %s", os.Args[1])
	}
}
