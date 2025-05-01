package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/storage"
)

func main() {
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// s := storage.NewNoStorage()
	s := storage.NewMemoryStorage()
	a := ai.NewAgent(ai.WithStorage(s))

	// wg := sync.WaitGroup{}

	// wg.Add(1)
	ai.StartRPCServer(a)

	// wg.Wait()
}
