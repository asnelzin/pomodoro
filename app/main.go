package main

import (
	"github.com/asnelzin/pomodoro/app/rest"
	"github.com/asnelzin/pomodoro/app/vk"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

var opts struct {
	SecretKey           string `long:"secret" env:"SECRET_KEY" description:"secret key"`
	APIToken            string `long:"token" env:"API_TOKEN" description:"vk api token"`
	ConfirmationMessage string `long:"confirmation" env:"CONFIRMATION_MESSAGE" description:"confirmation message"`
}

var revision string

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	log.Printf("pomodoro %s", revision)

	server := rest.Server{
		Bot:                 vk.NewBot(opts.APIToken),
		ConfirmationMessage: opts.ConfirmationMessage,
		SecretKey:           opts.SecretKey,
	}
	server.Run()
}
