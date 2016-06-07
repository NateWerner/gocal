package gcal

import (
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// GcalClient is a google calenar api client
type GcalClient struct {
	Srv  *calendar.Service
	Conf Config
}

// Event is google calendar event at Gcal
type Event struct {
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// NewCalendarClient returns  http client google calandar api
// scope is calendar.CalendarReadonlyScope or calendar.CalendarScope
func NewCalendarClient(c Config, scope string) (*GcalClient, error) {
	var gc GcalClient
	b, err := ioutil.ReadFile(c.Credential)
	if err != nil {
		return nil, err
	}

	jc, err := google.JWTConfigFromJSON(b, scope)
	if err != nil {
		return nil, err
	}

	//client := gcal.NewClient(ctx, config)
	client := jc.Client(oauth2.NoContext)

	srv, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	gc.Srv = srv
	gc.Conf = c
	return &gc, nil
}

// GetEventsList returns event list
func (gc GcalClient) GetEventsList(startTime string, endTime string) (*calendar.Events, error) {
	events, err := gc.Srv.Events.List(gc.Conf.CalendarID).TimeMax(endTime).
		TimeMin(startTime).SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	return events, nil
}

// InsertEvent insert an event to the google calendar
func (gc GcalClient) InsertEvent(event Event) error {
	start := calendar.EventDateTime{
		DateTime: event.StartTime,
	}
	end := calendar.EventDateTime{
		DateTime: event.EndTime,
	}

	ge := calendar.Event{
		Summary:     event.Title,
		Start:       &start,
		Description: event.Detail,
		End:         &end,
	}

	_, err := gc.Srv.Events.Insert(gc.Conf.CalendarID, &ge).Do()
	return err
}