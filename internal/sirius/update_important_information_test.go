package sirius

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateImportantInformation(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `{
    "complaints": "true",
    "panelDeputy":  "false",
    "annualBillingInvoice": "schedule",
    "otherImportantInformation": "This is some crucial information",
	}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	formData := ImportantInformationDetails{
		Complaints: true,
		PanelDeputy: false,
		AnnualBillingInvoice: "schedule",
		OtherImportantInformation: "This is some crucial information",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)
	assert.Nil(t, err)
}

func TestUpdateImportantInformationReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantInformationDetails{
		Complaints: true,
		PanelDeputy: false,
		AnnualBillingInvoice: "schedule",
		OtherImportantInformation: "This is some crucial information",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    fmt.Sprintf("%v/api/v1/deputies/%d/important-information", svr.URL, ID),
		Method: http.MethodPut,
	}, err)
}

func TestUpdateImportantInformationReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	formData := ImportantInformationDetails{
		Complaints: true,
		PanelDeputy: false,
		AnnualBillingInvoice: "schedule",
		OtherImportantInformation: "This is some crucial information",
	}

	err := client.UpdateImportantInformation(getContext(nil), ID, formData)

	assert.Equal(t, ErrUnauthorized, err)
}
