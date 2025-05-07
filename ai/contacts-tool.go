package ai

import (
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/pbreedt/ai-agent/contacts"
)

// var (
// 	contactsDB *contacts.SqliteContactsDB
// )

func GetContactTool(a *Agent) ai.Tool {

	getContactTool := genkit.LookupTool(a.genkit, "getContact")
	if getContactTool != nil {
		return getContactTool
	}

	getContactTool = genkit.DefineTool(
		a.genkit, "getContact", "Find and return a contact that matches the provided information.",
		func(ctx *ai.ToolContext, input contacts.Person) (string, error) {
			result, err := a.contactsDB.GetByPerson(input)
			if err != nil {
				return fmt.Sprintf("The contact could not be found due to the following error: %s", err.Error()), err
			}
			return fmt.Sprintf("The contact is: %s", result), nil
		})

	return getContactTool
}

func StoreContactTool(a *Agent) ai.Tool {

	storeContactTool := genkit.LookupTool(a.genkit, "storeContact")
	if storeContactTool != nil {
		return storeContactTool
	}

	storeContactTool = genkit.DefineTool(
		a.genkit, "storeContact", `Store the provided information of a contact. To store a contact, you need at least the full name (name and surname). 
		Any of the following details are optional: nickname, email, mobile, telegram id.
		If the contact already exists, it will be updated with the provided information.`,
		func(ctx *ai.ToolContext, input contacts.Person) (string, error) {
			p, e := a.contactsDB.GetByPerson(input)
			if e == nil && p != (contacts.Person{}) {
				if !input.Compare(p) {
					return "The contact was updated successfully.", a.contactsDB.Update(input)
				}
				return "The contact already exists: " + p.String(), nil
			}

			err := a.contactsDB.Insert(input)
			if err != nil {
				return fmt.Sprintf("The contact could not be stored due to the following error: %s", err.Error()), err
			}
			return "The contact was stored successfully.", nil
		})

	return storeContactTool
}
