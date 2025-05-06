package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/history"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s := history.NewInMemoryHistory(100) // TODO: make configurable
	a := ai.NewAgent(ai.WithHistory(s))

	ai.StartRPCServer(a)
}
