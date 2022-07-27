package main

import (
	"flag"
	tgClient "home/internal/clients/telegram"
	eventconsumer "home/internal/consumer/event-consumer"
	"home/internal/events/telegram"
	"home/pkg/lib/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "t", "", "Token for access to telegram bot.\nRequired for start bot.")
	flag.Parse()

	if token == "" {
		log.Fatal("token is not specified")
	}
}

func main() {
	eventsProcessor := telegram.New(tgClient.New(tgBotHost, token), files.New(storagePath))

	log.Println("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
