package main

import (
	"log"
	"os"

	tgClient "note-adviser-bot/clients/telegram"
	event_consumer "note-adviser-bot/consumer/event-consumer"
	"note-adviser-bot/events/telegram"
	"note-adviser-bot/storage/files"

	"github.com/joho/godotenv"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files-storage"
	batchSize   = 100
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token, exists := os.LookupEnv("NOTE_ADVISER_TOKEN")

	if !exists {
		log.Fatal("token is not specified")
	}

	return token
}
