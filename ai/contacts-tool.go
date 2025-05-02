package ai

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/pbreedt/ai-agent/contacts"
)

var (
	contactsDB *contacts.ContactsDB
)

func GetContactTool(g *genkit.Genkit) ai.Tool {

	getContactTool := genkit.LookupTool(g, "getContact")
	if getContactTool != nil {
		return getContactTool
	}

	if contactsDB == nil {
		contactsDB = contacts.NewContactsDB("../contacts.db")
		err := contactsDB.Open()
		if err != nil {
			panic(err)
		}
	}

	getContactTool = genkit.DefineTool(
		g, "getContact", "Find and return a contact that matches the provided information.",
		func(ctx *ai.ToolContext, input contacts.Person) (string, error) {
			// _, r := InitContactsIndexerRetriever(g)
			// rtrvFlow := Retrieve(g, r)
			// result, err := rtrvFlow.Run(context.Background(), fmt.Sprintf("Find the contact that matches the following details: %s", input.String()))

			result, err := contactsDB.GetByPerson(input)
			if err != nil {
				return err.Error(), err
			}
			return fmt.Sprintf("The contact is: %s", result), nil
		})

	return getContactTool
}
