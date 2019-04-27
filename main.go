package main

import (
	"log"
	"net/http"
)

func main() {
	config := Configure()

	OpenDB(&config)
	defer CloseDB()

	if config.Debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	log.Fatal(http.ListenAndServe(config.Addr, Router(Endpoints())))
}
