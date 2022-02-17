package main

import (
	"github.com/KarelKubat/smartlog/log"
)

func main() {
	log.Print("Hello world!")
	for i := 0; i < 10; i++ {
		log.Printf("Msg %d", i)
	}
	log.Fatalf("Bye %v.", "world")
}
