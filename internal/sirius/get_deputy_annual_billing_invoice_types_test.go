package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAnnualBillingInvoiceTypes(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
	{
            "handle": "INVOICE",
            "label": "Invoice"
        },
        {
            "handle": "SCHEDULE",
            "label": "Schedule"
        },
        {
            "handle": "SCHEDULE AND INVOICE",
            "label": "Schedule and Invoice"
        },
        {
            "handle": "UNKNOWN",
            "label": "Unknown"
        }
]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []DeputyAnnualBillingInvoiceTypes{
		{
			"INVOICE",
			"Invoice",
		},
		{
			"SCHEDULE",
			"Schedule",
		},
		{
			"SCHEDULE AND INVOICE",
			"Schedule and Invoice",
		},
		{
			"UNKNOWN",
			"Unknown",
		},
	}

	invoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(getContext(nil))

	assert.Equal(t, expectedResponse, invoiceTypes)
	assert.Equal(t, nil, err)
}

func TestGetAnnualBillingInvoiceTypesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	invoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(getContext(nil))

	assert.Equal(t, []DeputyAnnualBillingInvoiceTypes(nil), invoiceTypes)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reference-data/annualBillingInvoice",
		Method: http.MethodGet,
	}, err)
}

func TestGetAnnualBillingInvoiceTypesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	invoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(getContext(nil))

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, []DeputyAnnualBillingInvoiceTypes(nil), invoiceTypes)
}
