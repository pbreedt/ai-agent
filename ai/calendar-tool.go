package ai

import (
	"fmt"
	"log"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type CalendarInput struct {
	From time.Time `json:"from" doc:"Date and time to start looking for events in Calendar."`
	To   time.Time `json:"to" doc:"Date and time to stop looking for events in Calendar."`
	// From string `json:"from-datetime" doc:"Date and time in RFC3339 format. Calendar events that start after this date and time will be returned."`
	// To   string `json:"to-datetime,omitempty" doc:"Date and time in RFC3339 format. If not provided, then assume end date to be the same as the start date."`
}

func GetCalendarEventsTool(a *Agent) ai.Tool {

	getCalEvents := genkit.LookupTool(a.genkit, "getCalEvents")

	if getCalEvents != nil {
		return getCalEvents
	}

	getCalEvents = genkit.DefineTool(
		a.genkit, "getCalEvents", fmt.Sprintf(`
		Gets the calendar events from start date to end date. If no end date is provided, then assume end date to be the same as the start date.
		If a date is provided without a time, then assume the time is 00:00:00 for the start date and time is 23:59:59 for the end date.
		Date and time should be in RFC3339 format, i.e. "2022-01-01T00:00:00Z".
		Today's date is %s and the time now is %s. Using this date, you can infer other dates, such as next week or last month.`, time.Now().Format(time.DateOnly), time.Now().Format(time.TimeOnly)),
		func(ctx *ai.ToolContext, input CalendarInput) (string, error) {
			if input.From.IsZero() {
				return "Please provide at least a start date.", nil
			}
			// if input.To.IsZero() {
			// 	// Set input.to to be the same as input.from with time 23:59:59
			// 	input.To = time.Date(input.From.Year(), input.From.Month(), input.From.Day(), 23, 59, 59, 0, time.UTC)
			// }

			log.Printf("[GetCalendarEventsTool] Getting events from %s to %s\n", input.From, input.To)

			events, err := a.calendar.GetEvents(input.From, input.To)
			if err != nil {
				return fmt.Sprintf("The calendar entries could not be found due to the following error: %s", err.Error()), err
			}

			var eventStrings []string
			for _, event := range events {
				eventStrings = append(eventStrings, event.String())
			}
			return fmt.Sprintf("The calendar entries from %s to %s are: %s", input.From, input.To, eventStrings), nil
		})

	return getCalEvents
}
