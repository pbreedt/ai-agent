package calendar

import (
	"fmt"
	"time"
)

type Event struct {
	Summary   string
	Start     string
	End       string
	Location  string
	Attendees []*EventAttendee
}

type EventAttendee struct {
	Email       string
	DisplayName string
}

func (e Event) String() string {
	return fmt.Sprintf("Event: %s\nStart: %s\nEnd: %s\nLocation: %s\nAttendees: %s", e.Summary, e.Start, e.End, e.Location, e.Attendees)
}

func (e EventAttendee) String() string {
	return e.DisplayName
}

type Calendar interface {
	GetEvents(from time.Time, to time.Time) ([]Event, error)
	CreateEvent(event *Event) (*Event, error)
}
