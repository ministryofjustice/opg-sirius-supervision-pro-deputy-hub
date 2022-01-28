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

func TestChangeECM(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{"ecmId": 32}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, ProDeputyDetails{ID: 76})
	assert.Equal(t, nil, err)
}

func TestChangeECMReturnsErrorIfNoEcm(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 0}

	err := client.ChangeECM(getContext(nil), changeEcmForm, ProDeputyDetails{ID: 76})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/ecm",
		Method: http.MethodPut,
	}, err)
}

func TestChangeECMReturnsErrorIfNoId(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, ProDeputyDetails{ID: 0})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/0/ecm",
		Method: http.MethodPut,
	}, err)
}

func TestChangeECMReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, ProDeputyDetails{ID: 76})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/76/ecm",
		Method: http.MethodPut,
	}, err)
}

func TestChangeECMReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)
	changeEcmForm := ExecutiveCaseManagerOutgoing{EcmId: 23}

	err := client.ChangeECM(getContext(nil), changeEcmForm, ProDeputyDetails{ID: 76})

	assert.Equal(t, ErrUnauthorized, err)
}
