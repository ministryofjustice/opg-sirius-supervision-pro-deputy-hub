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

func TestProDeputyDetailsReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `    {
		"id": 76,
		"firstname": "firstname",
		"surname": "surname",
		"deputyNumber": 1000,
		"deputyStatus": "INACTIVE",
		"organisationName": "organisationName",
		"executiveCaseManager": {
			"id": 223,
    		"displayName": "displayName"
		},
		"firm": {
			"firmName": "This is the Firm Name"
		}
    }`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := ProDeputyDetails{
		ID:               76,
		DeputyFirstName:  "firstname",
		DeputySurname:    "surname",
		DeputyNumber:     1000,
		DeputyStatus:     "INACTIVE",
		OrganisationName: "organisationName",
		ExecutiveCaseManager: executiveCaseManager{
			EcmId:   223,
			EcmName: "displayName",
		},
		Firm: firm{
			FirmName: "This is the Firm Name",
		},
	}

	proDeputyDetails, err := client.GetProDeputyDetails(getContext(nil), 76)

	assert.Equal(t, expectedResponse, proDeputyDetails)
	assert.Equal(t, nil, err)
}

func TestGetDeputyDetailsReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetProDeputyDetails(getContext(nil), 76)

	expectedResponse := ProDeputyDetails{
		ID:               0,
		DeputyFirstName:  "",
		DeputySurname:    "",
		DeputyNumber:     0,
		DeputyStatus:     "",
		OrganisationName: "",
		ExecutiveCaseManager: executiveCaseManager{
			EcmId:   0,
			EcmName: "",
		},
	}

	assert.Equal(t, expectedResponse, deputyDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyDetailsReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	deputyDetails, err := client.GetProDeputyDetails(getContext(nil), 76)

	expectedResponse := ProDeputyDetails{
		ID:               0,
		DeputyFirstName:  "",
		DeputySurname:    "",
		DeputyNumber:     0,
		DeputyStatus:     "",
		OrganisationName: "",
		ExecutiveCaseManager: executiveCaseManager{
			EcmId:   0,
			EcmName: "",
		},
	}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, deputyDetails)
}
