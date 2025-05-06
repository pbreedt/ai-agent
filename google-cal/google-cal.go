package googlecal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pbreedt/ai-agent/ai"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct {
	srv        *calendar.Service
	maxResults int
}

func NewGoogleCalendar(maxResults int) (*GoogleCalendar, error) {
	ctx := context.Background()
	b, err := os.ReadFile("google-cal-creds.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarEventsScope, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return nil, err
	}

	return &GoogleCalendar{srv: srv, maxResults: maxResults}, nil

	// ShowAllCalendars(srv)

	// t := time.Now().Format(time.RFC3339)

	// events, err := srv.Events.List("primary").ShowDeleted(false).
	// 	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	// }
	// fmt.Println("Upcoming events:")
	// if len(events.Items) == 0 {
	// 	fmt.Println("No upcoming events found.")
	// } else {
	// 	for _, item := range events.Items {
	// 		date := item.Start.DateTime
	// 		if date == "" {
	// 			date = item.Start.Date
	// 		}
	// 		fmt.Printf("%v (%v)\n", item.Summary, date)
	// 	}
	// }
}

func (cal *GoogleCalendar) GetEvents(from time.Time, to time.Time) ([]ai.CalendarEvent, error) {
	f := from.Format(time.RFC3339)
	t := to.Format(time.RFC3339)

	log.Printf("Getting events from %s to %s\n", f, t)

	events, err := cal.srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(f).TimeMax(t).MaxResults(int64(cal.maxResults)).OrderBy("startTime").Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		return nil, err
	}

	var r []ai.CalendarEvent
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
			r = append(r, ai.CalendarEvent{
				Summary:  item.Summary,
				Start:    date,
				End:      item.End.DateTime,
				Location: item.Location,
				// Attendees: item.Attendees, TODO: add attendees
			})
		}
	}

	return r, nil
}

func GetEventsForDate(srv *calendar.Service, date time.Time) {
	t := date.Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).TimeMax(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	fmt.Println("Upcoming events for", t, ":")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}

// start web server and listen for requests
func ShowAllCalendars(srv *calendar.Service) {
	fmt.Println("Calendars:")

	calList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve calendar list: %v", err)
	}

	for _, cal := range calList.Items {
		fmt.Println(cal.Id, cal.Summary)
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Login using link below:\n%v\n", authURL)

	fmt.Print("Enter authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
