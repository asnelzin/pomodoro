package main

import (
	"github.com/asnelzin/pomodoro/app/rest"
	"log"
)

var revision string

func main() {
	log.Printf("pomodoro %s", revision)

	server := rest.Server{}
	server.Run()
}