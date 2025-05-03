package contacts

import (
	"encoding/json"
	"fmt"
	"log"
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
	if e != nil {
		log.Println("JSON marshal error:", e.Error())

		p := "{"
		if c.Id != nil {
			p += fmt.Sprintf(`"db_id":%d`, *(c.Id))
		}
		if c.Name != "" {
			p += fmt.Sprintf(`,"name":"%s"`, c.Name)
		}
		if c.Surname != "" {
			p += fmt.Sprintf(`,"surname":"%s"`, c.Surname)
		}
		if c.Nickname != "" {
			p += fmt.Sprintf(`,"nickname":"%s"`, c.Nickname)
		}
		if c.Email != "" {
			p += fmt.Sprintf(`,"email":"%s"`, c.Email)
		}
		if c.Mobile != "" {
			p += fmt.Sprintf(`,"mobile":"%s"`, c.Mobile)
		}
		if c.TelegramID != nil {
			p += fmt.Sprintf(`,"telegram-id":%d`, *(c.TelegramID))
		}

		return p + "}"
	}
	return string(b)
}

func (c Person) ToMap() map[string]any {
	return map[string]any{
		"db_id":       *c.Id,
		"name":        c.Name,
		"surname":     c.Surname,
		"nickname":    c.Nickname,
		"email":       c.Email,
		"mobile":      c.Mobile,
		"telegram-id": *c.TelegramID,
	}
}

func (c Person) Compare(other Person) bool {
	match := true
	if c.Id != nil && other.Id != nil {
		match = *(c.Id) == *(other.Id)
	}
	if c.Id != nil && other.Id == nil {
		match = false
	}
	if c.Id == nil && other.Id != nil {
		match = false
	}
	match = match && c.Name == other.Name
	match = match && c.Surname == other.Surname
	match = match && c.Nickname == other.Nickname
	match = match && c.Email == other.Email
	match = match && c.Mobile == other.Mobile
	if c.TelegramID != nil && other.TelegramID != nil {
		match = *(c.TelegramID) == *(other.TelegramID)
	}
	if c.TelegramID != nil && other.TelegramID == nil {
		match = false
	}
	if c.TelegramID == nil && other.TelegramID != nil {
		match = false
	}
	return match
}
