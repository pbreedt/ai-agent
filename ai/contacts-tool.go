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

func StoreContactTool(g *genkit.Genkit) ai.Tool {

	storeContactTool := genkit.LookupTool(g, "storeContact")
	if storeContactTool != nil {
		return storeContactTool
	}

	if contactsDB == nil {
		contactsDB = contacts.NewContactsDB("../contacts.db")
		err := contactsDB.Open()
		if err != nil {
			panic(err)
		}
	}

	storeContactTool = genkit.DefineTool(
		g, "storeContact", `Store the provided information of a contact. To store a contact, you need at least the full name (name and surname). 
		Any of the following details are optional: nickname, email, mobile, telegram id.
		If the contact already exists, it will be updated with the provided information.`,
		func(ctx *ai.ToolContext, input contacts.Person) (string, error) {
			p, e := contactsDB.GetByPerson(input)
			if e == nil && p != (contacts.Person{}) {
				if !input.Compare(p) {
					return "The contact was updated successfully.", contactsDB.Update(input)
				}
				return "The contact already exists: " + p.String(), nil
			}

			err := contactsDB.Insert(input)
			if err != nil {
				return err.Error(), err
			}
			return "The contact was stored successfully.", nil
		})

	return storeContactTool
}

// Opted to handle update and insert together in 'storeContact' - more control over duplicate records
// func UpdateContactTool(g *genkit.Genkit) ai.Tool {

// 	updateContactTool := genkit.LookupTool(g, "updateContact")
// 	if updateContactTool != nil {
// 		return updateContactTool
// 	}

// 	if contactsDB == nil {
// 		contactsDB = contacts.NewContactsDB("../contacts.db")
// 		err := contactsDB.Open()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	updateContactTool = genkit.DefineTool(
// 		g, "updateContact", "Update the provided information of a contact. To update a contact, you need either the id the full name (name and surname) of the contact. Any other provided details (nickname, email, mobile, telegram id) should be updated.",
// 		func(ctx *ai.ToolContext, input contacts.Person) (string, error) {
// 			err := contactsDB.Update(input)
// 			if err != nil {
// 				return err.Error(), err
// 			}
// 			return "The contact was stored successfully.", nil
// 		})

// 	return updateContactTool
// }
