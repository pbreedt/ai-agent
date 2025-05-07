package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pbreedt/ai-agent/contacts"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// ctx := context.Background()

	// Used Google AI instead of Vertex AI
	// g, err := genkit.Init(ctx,
	// 	genkit.WithPlugins(&googlegenai.GoogleAI{}),
	// 	genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	db := contacts.NewSqliteContactsDB("../contacts.db")
	e := db.Open()
	if e != nil {
		log.Fatalf("could not open database: %v", e)
	}
	defer db.Close()
	e = db.CreateContactTable()
	if e != nil {
		log.Fatalf("could not create contact table: %v", e)
	}

	// p := contacts.Person{
	// 	Name:     "Petrus",
	// 	Surname:  "Breedt",
	// 	Nickname: "Peet",
	// 	Email:    "petrus.breedt@gmail.com",
	// 	Mobile:   "+1 (224) 706-7025",
	// }

	// e = db.Insert(p)
	// if e != nil {
	// 	log.Fatalf("could not insert contact: %v", e)
	// }

	// p = contacts.Person{
	// 	Name:     "John",
	// 	Surname:  "Doe",
	// 	Nickname: "John Doe",
	// 	Email:    "jdoe@example.com",
	// 	Mobile:   "+31612345678",
	// }

	// e = db.Insert(p)
	// if e != nil {
	// 	log.Fatalf("could not insert contact: %v", e)
	// }

	all, err := db.GetAll()
	if err != nil {
		log.Printf("error in retrieving all: %v\n", err)
	}
	for _, c := range all {
		log.Printf("%+v\n", c)
	}

	// i, _ := ai.InitContactsIndexerRetriever(g)
	// idxFlow := ai.Index(g, i)
	// _, e := idxFlow.Run(context.Background(), c)
	// if e != nil {
	// 	fmt.Printf("error in indexing: %v\n", e)
	// }

	// rtrvFlow := ai.Retrieve(g, r)
	// result, err := rtrvFlow.Run(context.Background(), fmt.Sprintf("Find the contact that matches the following details: %s", c.String()))
	// if err != nil {
	// 	fmt.Printf("error in retrieving: %v\n", err)
	// }
	// fmt.Printf("result: %s\n", result)

}
