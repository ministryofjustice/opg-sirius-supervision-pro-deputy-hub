package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeputyEventsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
    {
      "id": 300,
      "hash": "AW",
      "timestamp": "2021-09-09 14:01:59",
      "eventType": "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder",
      "user": {
        "id": 41,
        "phoneNumber": "12345678",
        "displayName": "system admin",
        "email": "system.admin@opgtest.com"
      },
      "event": {
        "orderType": "pfa",
        "orderUid": "7000-0000-1995",
        "orderId": "58",
        "orderCourtRef": "03305972",
        "courtReferenceNumber": "03305972",
        "courtReference": "03305972",
        "personType": "Deputy",
        "personId": "76",
        "personUid": "7000-0000-2530",
        "personName": "Mx Bob Builder",
        "personCourtRef": null,
        "additionalPersons": [
          {
            "personType": "Client",
            "personId": "63",
            "personUid": "7000-0000-1961",
            "personName": "Test Name",
            "personCourtRef": "40124126"
          }
        ]
      }
    }
  ]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := ProDeputyEventCollection{
		ProDeputyEvent{
			TimelineEventId: 300,
			Timestamp:       "2021-09-09 14:01:59",
			EventType:       "DeputyLinkedToOrder",
			User:            User{UserId: 41, UserDisplayName: "system admin", UserPhoneNumber: "12345678"},
			Event: Event{
				DeputyID:    "76",
				DeputyName:  "Mx Bob Builder",
				OrderType:   "pfa",
				SiriusId:    "7000-0000-1995",
				OrderNumber: "03305972",
				Client:      []ClientPerson{{ClientName: "Test Name", ClientId: "63", ClientUid: "7000-0000-1961", ClientCourtRef: "40124126"}},
			},
		},
	}

	deputyEvents, err := client.GetProDeputyTimeline(getContext(nil), 1)

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, nil, err)
}

func TestGetDeputyEventsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetProDeputyTimeline(getContext(nil), 76)

	expectedResponse := ProDeputyEventCollection(nil)

	assert.Equal(t, expectedResponse, deputyEvents)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/timeline/76",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyEventsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyEvents, err := client.GetProDeputyTimeline(getContext(nil), 76)

	expectedResponse := ProDeputyEventCollection(nil)

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyEvents)
}

func TestFormatDateAndTime(t *testing.T) {
	unsortedData := "2020-10-18 10:11:08"
	expectedResponse := "18/10/2020 10:11:08"
	assert.Equal(t, expectedResponse, reformatTimestamp(unsortedData))
}

func TestReformatEventType(t *testing.T) {
	expectedResponse := "DeputyLinkedToOrder"
	testDeputyEvent := "Opg\\Core\\Model\\Event\\Order\\DeputyLinkedToOrder"
	assert.Equal(t, expectedResponse, reformatEventType(testDeputyEvent))
}
