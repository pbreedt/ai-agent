package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/pbreedt/ai-agent/ai"
)

func main() {
	aic, err := ai.GetRPCClient()
	if err != nil {
		log.Fatal("Error getting RPC client:", err)
	}
	defer aic.Close()

	var aiRes string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Hi there. How can I help you?\n\n> ")
	for {
		scanner.Scan()
		prompt := scanner.Text()
		switch prompt {
		case "exit":
			fmt.Println("OK, bye!")
			return
		case "help":
			fmt.Print("Type a question or instruction and I will respond.\nAlternatively, type 'exit' to quit.\n\n> ")
			continue
		case "":
			fmt.Print("Awaiting your next question or instructions (or type 'help').\n\n> ")
			continue
		}

		err = aic.Call("Agent.RPCRespondToPrompt", prompt, &aiRes)
		if err != nil {
			aiRes = err.Error()
		}

		fmt.Print(aiRes, "\n> ")
		aiRes = ""

		if aiRes == "AI shut down" {
			return
		}
	}

}
