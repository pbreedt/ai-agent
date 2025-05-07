package ai

import (
	"github.com/pbreedt/ai-agent/calendar"
	"github.com/pbreedt/ai-agent/contacts"
	"github.com/pbreedt/ai-agent/history"
)

type Option func(a *Agent)

// WithChatHistory sets the storage to be used for keeping chat history
func WithChatHistory(history history.ChatHistory) Option {
	return func(a *Agent) { a.chatHistory = history }
}

func WithContactsDB(db contacts.ContactsDB) Option {
	return func(a *Agent) { a.contactsDB = db }
}

// WithCalendar sets the calendar service to be used for getting/storing calendar events
func WithCalendar(cal calendar.Calendar) Option {
	return func(a *Agent) { a.calendar = cal }
}
