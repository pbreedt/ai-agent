package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/calendar"
	"github.com/pbreedt/ai-agent/contacts"
	"github.com/pbreedt/ai-agent/history"
)

// run with 'genkit start -- go run .' to access Genkit UI with under-the-hood insights
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up calendar
	gcal, err := calendar.NewGoogleCalendar(10)
	if err != nil {
		log.Printf("Error creating Google Calendar client: %v", err)
	}
	// gcal.GetEvents(time.Now(), time.Now().Add(time.Hour*24))

	// Set up chat history
	s := history.NewInMemoryHistory(100) // TODO: make configurable

	// Set up contacts DB
	conDb := contacts.NewSqliteContactsDB("../contacts.db")
	err = conDb.Open()
	if err != nil {
		log.Printf("Error opening contacts DB: %v", err)
	}

	a := ai.NewAgent(
		ai.WithChatHistory(s),
		ai.WithCalendar(gcal),
		ai.WithContactsDB(conDb),
	)

	ai.StartRPCServer(a)
}
