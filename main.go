package main

import (
	"log"
)

func main() {
	s := NewServer()

	s.Start()
	defer s.Stop()

	if s.Config.Debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}
