package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Debug      bool   `default:"false"`
	Addr       string `default:":8080"`
	DBHost     string `default:"postgres"`
	DBPort     string `default:"5432"`
	DBUser     string `default:"postgres"`
	DBPassword string `default:""`
	DBName     string `default:"postgres"`
	DBSslMode  string `default:"disable"`
}

func Configure() Configuration {
	var config Configuration

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}
