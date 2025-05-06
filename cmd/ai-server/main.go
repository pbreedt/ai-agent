package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	googlecal "github.com/pbreedt/ai-agent/google-cal"
	"github.com/pbreedt/ai-agent/history"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gcal, err := googlecal.NewGoogleCalendar(10)
	if err != nil {
		log.Printf("Error creating Google Calendar client: %v", err)
	}
	// gcal.GetEvents(time.Now(), time.Now().Add(time.Hour*24))

	s := history.NewInMemoryHistory(100) // TODO: make configurable
	a := ai.NewAgent(ai.WithChatHistory(s), ai.WithCalendar(gcal))

	ai.StartRPCServer(a)
}
