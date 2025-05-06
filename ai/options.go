package ai

import (
	"fmt"
	"time"

	"github.com/firebase/genkit/go/ai"
)

type Option func(b *Agent)

type ChatHistory interface {
	GetLast(last int) []*ai.Message
	GetAll() []*ai.Message
	Store(message *ai.Message) error
}

// WithChatHistory sets the storage to be used for keeping chat history
func WithChatHistory(history ChatHistory) Option {
	return func(b *Agent) { b.chatHistory = history }
}

type CalendarEvent struct {
	Summary   string
	Start     string
	End       string
	Location  string
	Attendees []string
}

func (ce CalendarEvent) String() string {
	return fmt.Sprintf("Event: %s\nStart: %s\nEnd: %s\nLocation: %s\nAttendees: %s", ce.Summary, ce.Start, ce.End, ce.Location, ce.Attendees)
}

type Calendar interface {
	GetEvents(from time.Time, to time.Time) ([]CalendarEvent, error)
}

func WithCalendar(cal Calendar) Option {
	return func(b *Agent) { b.calendar = cal }
}
