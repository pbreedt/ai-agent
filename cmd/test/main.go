package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/pbreedt/ai-agent/ai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	g, _, r := ai.InitContacts()
	// idxFlow := ai.Index(g, i)
	// idxFlow.Run(context.Background(), ai.Contact{
	// 	Name:     "John",
	// 	Surname:  "Doe",
	// 	Nickname: "John Doe",
	// 	Email:    "jdoe@example.com",
	// 	Mobile:   "+31612345678"})

	rtrvFlow := ai.Retrieve(g, r)
	result, err := rtrvFlow.Run(context.Background(), "Jim Doe")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("result: %s\n", result)

}
