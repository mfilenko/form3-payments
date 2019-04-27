package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config := Configure()
	OpenDB(&config)
	defer CloseDB()

	os.Exit(m.Run())
}
