package main

import (
	"github.com/asnelzin/pomodoro/app/rest"
	"log"
)

var revision string

func main() {
	log.Printf("secrets %s", revision)

	server := rest.Server{}
	server.Run()
}