package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s := storage.NewMemoryStorage()
	a := ai.NewAgent(ai.WithStorage(s))

	ai.StartRPCServer(a)
}
