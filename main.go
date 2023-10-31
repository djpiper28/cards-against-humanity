package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC)
	log.Println("Starting up Cards Against Humanity server")
}
