package main

import (
	"github.com/asnelzin/pomodoro/app/rest"
	"log"
	"os"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	SecretKey string `long:"secret" env:"SECRET_KEY" description:"secret key"`
	APIToken string `long:"token" env:"API_TOKEN" description:"vk api token"`
	ConfirmationMessage string `long:"confirmation" env:"CONFIRMATION_MESSAGE" description:"confirmation message"`
}

var revision string

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	//bot := getBot(opts.APIToken)

	log.Printf("pomodoro %s", revision)

	server := rest.Server{
		ConfirmationMessage: opts.ConfirmationMessage,
		SecretKey: opts.SecretKey,
	}
	server.Run()
}