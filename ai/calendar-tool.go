package ai

import (
	"fmt"
	"log"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/pbreedt/ai-agent/calendar"
)

type CalendarInput struct {
	From time.Time `json:"from" doc:"Date and time to start looking for events in Calendar."`
	To   time.Time `json:"to" doc:"Date and time to stop looking for events in Calendar."`
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
		Date and time should be in RFC3339 format, i.e. "2022-01-01T00:00:00-05:00" for Chicago time.
		Today's date is %s and the time now is %s.
		Using this date, you can infer other dates, such as next week or last month which will have start and end dates relative to today.
		The week starts at 00:00:00 on Monday and ends at 23:59:59 on Sunday. When using relative dates, also confirm which dates you're using.
		`, time.Now().Format(time.DateOnly), time.Now().Format(time.TimeOnly)),
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

func CreateCalendarEventsTool(a *Agent) ai.Tool {

	createCalEvents := genkit.LookupTool(a.genkit, "createCalEvents")

	if createCalEvents != nil {
		return createCalEvents
	}

	createCalEvents = genkit.DefineTool(
		a.genkit, "createCalEvents", fmt.Sprintf(`
		Create and save a calendar event. If no end date is provided, then assume end date to be the same as the start date.
		The following information can be provided:
		- summary
		- start date and time
		- end date and time
		- location
		- attendees
		At least a sammary and start date and time must be provided. 
		If no end date is provided, then assume end date to be the same as the start date.
		If no end time is provided, then assume end time to be 1 hour after the start time.
		Date and time should be in RFC3339 format, i.e. "2022-01-01T00:00:00-05:00" for Chicago time.
		Today's date is %s and the time now is %s.
		Using this date, you can infer other dates, such as next week or last month which will have start and end dates relative to today.
		The week starts at 00:00:00 on Monday and ends at 23:59:59 on Sunday. When using relative dates, also confirm which dates you're using.
		If the event was created successfully, respond with a summary of the event details.
		`, time.Now().Format(time.DateOnly), time.Now().Format(time.TimeOnly)),
		func(ctx *ai.ToolContext, input calendar.Event) (string, error) {
			if len(input.Start) == 0 {
				return "Please provide at least a start date.", nil
			}

			log.Printf("[CreateCalendarEventsTool] Creating event: %s\n", input)

			stored, err := a.calendar.CreateEvent(&input)
			if err != nil {
				return fmt.Sprintf("The event could not be created due to the following error: %s", err.Error()), err
			}

			return fmt.Sprintf("The calendar event was successfully created: %s", stored), nil
		})

	return createCalEvents
}
