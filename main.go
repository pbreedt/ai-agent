package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"github.com/pbreedt/ai-agent/ai"
	"github.com/pbreedt/ai-agent/storage"
	"github.com/pbreedt/ai-agent/telegram"
)

// Env:
// export TELEGRAM_BOT_TOKEN=<token>
// export GEMINI_API_KEY=<key>
// Setup:
// - AI Genkit
// - Memory Storage
// - Telegram Bot

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// s := storage.NewNoStorage()
	s := storage.NewMemoryStorage()
	a := ai.NewAgent(ai.WithStorage(s))

	t := telegram.New(a)

	go cli(ctx, a)

	t.StartListner(ctx)
}

// Keep reading line by line from std in and respond to std out
func cli(ctx context.Context, a *ai.Agent) {

	promptChan := make(chan string, 1)
	respChan := make(chan string)

	go a.Run(ctx, promptChan, respChan)

	defer close(promptChan)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Hi there. How can I help you?\n\n> ")

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println("Context done (main):", ctx.Err())
			return
		default:
		}

		switch scanner.Text() {
		case "exit":
			fmt.Println("OK, bye!\n")
			return
		case "help":
			fmt.Print("Type a question or instruction and I will respond.\nAlternatively, type 'exit' to quit.\n\n> ")
			continue
		case "":
			fmt.Print("Awaiting your next question or instructions (or type 'help').\n\n> ")
			continue
		}

		select {
		case <-ctx.Done():
			log.Println("Context done (main):", ctx.Err())
			return
		case promptChan <- scanner.Text():
		default:
			log.Println("Prompt channel full")
			continue
		}
		// log.Println("Sent prompt to AI")

		select {
		case <-ctx.Done():
			log.Println("Context done (main):", ctx.Err())
			return
		case res, ok := <-respChan:
			if !ok || res == "AI shut down" {
				log.Println("CLI shut down")
				return
			}
			fmt.Print(res, "\n> ")
		}
	}
}
