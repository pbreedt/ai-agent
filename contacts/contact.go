package contacts

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Id       *int   `json:"db_id,omitempty" doc:"Id of the contact in the database"`
	Name     string `json:"name,omitempty" doc:"contacts first name"`
	Surname  string `json:"surname,omitempty" doc:"contact's last name"`
	Nickname string `json:"nickname,omitempty" doc:"contact's nickname, alias or preferred name"`
	// Relationship string `json:"relationship,omitempty" doc:"contact's role or relationship to the user"`
	Email      string `json:"email,omitempty" doc:"contact's email address"`
	Mobile     string `json:"mobile,omitempty" doc:"contact's mobile number"`
	TelegramID *int64 `json:"telegram-id,omitempty" doc:"contact's telegram user (chat) id"`
}

func IntPointer(i int) *int {
	return &i
}
func Int64Pointer(i int64) *int64 {
	return &i
}

// return the contact in JSON format using json.Marshal
func (c Person) String() string {
	b, e := json.Marshal(c)
	if e == nil {
		return fmt.Sprintf(`{
			"db_id": %d,
			"name": "%s",
			"surname": "%s",
			"nickname": "%s",
			"email": "%s",
			"mobile": "%s",
			"telegram-id": "%d"
		}`, *c.Id, c.Name, c.Surname, c.Nickname, c.Email, c.Mobile, c.TelegramID)
	}
	return string(b)
}

func (c Person) ToMap() map[string]any {
	return map[string]any{
		"db_id":       c.Id,
		"name":        c.Name,
		"surname":     c.Surname,
		"nickname":    c.Nickname,
		"email":       c.Email,
		"mobile":      c.Mobile,
		"telegram-id": c.TelegramID,
	}
}
