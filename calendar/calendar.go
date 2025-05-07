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
	Attendees []string
}

func (e Event) String() string {
	return fmt.Sprintf("Event: %s\nStart: %s\nEnd: %s\nLocation: %s\nAttendees: %s", e.Summary, e.Start, e.End, e.Location, e.Attendees)
}

type Calendar interface {
	GetEvents(from time.Time, to time.Time) ([]Event, error)
}
