package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ProDeputyEventCollection []ProDeputyEvent

type User struct {
	UserId          int    `json:"id"`
	UserDisplayName string `json:"displayName"`
	UserPhoneNumber string `json:"phoneNumber"`
}

type Event struct {
	OrderType        string         `json:"orderType"`
	SiriusId         string         `json:"orderUid"`
	OrderNumber      string         `json:"orderCourtRef"`
	DeputyID         string         `json:"personId"`
	DeputyName       string         `json:"personName"`
	OrganisationName string         `json:"organisationName"`
	Changes          []Changes      `json:"changes"`
	Client           []ClientPerson `json:"additionalPersons"`
}

type Changes struct {
	FieldName string `json:"fieldName"`
	OldValue  string `json:"oldValue"`
	NewValue  string `json:"newValue"`
}

type ClientPerson struct {
	ClientName     string `json:"personName"`
	ClientId       string `json:"personId"`
	ClientUid      string `json:"personUid"`
	ClientCourtRef string `json:"personCourtRef"`
}

type ProDeputyEvent struct {
	TimelineEventId int    `json:"id"`
	Timestamp       string `json:"timestamp"`
	EventType       string `json:"eventType"`
	User            User   `json:"user"`
	Event           Event  `json:"event"`
}

func (c *Client) GetProDeputyTimeline(ctx Context, deputyId int) (ProDeputyEventCollection, error) {
	var v ProDeputyEventCollection

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/timeline/%d", deputyId), nil)

	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&v)

	DeputyEvents := editProDeputyEvents(v)

	return DeputyEvents, err

}

func editProDeputyEvents(v ProDeputyEventCollection) ProDeputyEventCollection {
	var list ProDeputyEventCollection
	for _, s := range v {
		event := ProDeputyEvent{
			Timestamp:       reformatTimestamp(s.Timestamp),
			EventType:       reformatEventType(s.EventType),
			TimelineEventId: s.TimelineEventId,
			User:            s.User,
			Event:           s.Event,
		}

		list = append(list, event)
	}
	return sortByTimestampAsc(list)
}

func sortByTimestampAsc(v ProDeputyEventCollection) ProDeputyEventCollection {
	sort.Slice(v, func(i, j int) bool {
		changeToTimeTypeI, _ := time.Parse("02/01/2006 15:04:05", v[i].Timestamp)
		changeToTimeTypeJ, _ := time.Parse("02/01/2006 15:04:05", v[j].Timestamp)
		return changeToTimeTypeJ.Before(changeToTimeTypeI)
	})
	return v
}

func reformatTimestamp(dateString string) string {
	dateTime, _ := time.Parse("2006-01-02 15:04:05", dateString)
	return dateTime.Format("02/01/2006 15:04:05")
}

func reformatEventType(event string) string {
	eventTypeArray := strings.Split(event, "\\")
	eventTypeArrayLength := len(eventTypeArray)
	eventTypeName := eventTypeArray[eventTypeArrayLength-1]
	return eventTypeName
}
