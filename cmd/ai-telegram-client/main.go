package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	aic, err := ai.GetRPCClient()
	if err != nil {
		log.Fatal("Error getting RPC client:", err)
	}
	defer aic.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	telegram.StartListner(ctx, aic)
}
