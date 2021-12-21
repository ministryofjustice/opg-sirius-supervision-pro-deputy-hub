package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEditPii(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
		"piiReceived":"20/01/2020",
		"piiExpiry":"20/01/2025",
		"piiAmount":"254",
		"piiRequested":"10/01/2020",
		}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 201,
			Body:       r,
		}, nil
	}

	piiDetails := PiiDetails{
		PiiReceived:  "20/01/2020",
		PiiExpiry:    "20/01/2025",
		PiiAmount:    "254",
		PiiRequested: "10/01/2020",
	}

	err := client.EditPiiCertificate(getContext(nil), piiDetails)
	assert.Nil(t, err)
}

// func TestAddFirmReturnsNewStatusError(t *testing.T) {
// 	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}))
// 	defer svr.Close()

// 	client, _ := NewClient(http.DefaultClient, svr.URL)

// 	int, err := client.AddFirmDetails(getContext(nil), FirmDetails{})

// 	assert.Equal(t, StatusError{
// 		Code:   http.StatusMethodNotAllowed,
// 		URL:    svr.URL + "/api/v1/firms",
// 		Method: http.MethodPost,
// 	}, err)

// 	assert.Equal(t, 0, int)
// }

// func TestAddFirmReturnsUnauthorisedClientError(t *testing.T) {
// 	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 	}))
// 	defer svr.Close()

// 	client, _ := NewClient(http.DefaultClient, svr.URL)

// 	int, err := client.AddFirmDetails(getContext(nil), FirmDetails{})

// 	assert.Equal(t, ErrUnauthorized, err)
// 	assert.Equal(t, 0, int)

// }
