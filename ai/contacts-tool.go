package ai

import (
	"encoding/json"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type Contact struct {
	Name       string `json:"name,omitempty" doc:"contacts first name"`
	Surname    string `json:"surname,omitempty" doc:"contact's last name"`
	Nickname   string `json:"nickname,omitempty" doc:"contact's nickname, alias or preferred name"`
	Email      string `json:"email,omitempty" doc:"contact's email address"`
	Mobile     string `json:"mobile,omitempty" doc:"contact's mobile number"`
	TelegramID string `json:"telegram-id,omitempty" doc:"contact's telegram user (chat) id"`
}

// return the contact in JSON format using json.Marshal
func (c Contact) String() string {
	b, e := json.Marshal(c)
	if e == nil {
		return fmt.Sprintf(`{
			"name": "%s",
			"surname": "%s",
			"nickname": "%s",
			"email": "%s",
			"mobile": "%s",
			"telegram-id": "%s"
		}`, c.Name, c.Surname, c.Nickname, c.Email, c.Mobile, c.TelegramID)
	}
	return string(b)
}

func (c Contact) ToMap() map[string]any {
	return map[string]any{
		"name":        c.Name,
		"surname":     c.Surname,
		"nickname":    c.Nickname,
		"email":       c.Email,
		"mobile":      c.Mobile,
		"telegram-id": c.TelegramID,
	}
}

func GetContactTool(g *genkit.Genkit) ai.Tool {

	getContactTool := genkit.LookupTool(g, "getContact")

	if getContactTool != nil {
		return getContactTool
	}

	getContactTool = genkit.DefineTool(
		g, "getContact", "Find and return a contact that matches the provided information.",
		func(ctx *ai.ToolContext, input Contact) (string, error) {
			// g2, _, r := ai.InitContacts()
			// Here, we would typically make an API call or database query.
			// For this example, we just return a fixed value.
			return fmt.Sprintf("Find a contact that matches the provided information: %s. Return all the information you have about this contact. If you cannot find the contact, say 'Contact not found'. Do not make up an answer.", input), nil
		})

	return getContactTool
}
