package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockManagePiiDetailsInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.ProDeputyDetails
}

func (m *mockManagePiiDetailsInformation) EditPiiCertificate(ctx sirius.Context, piiData sirius.PiiDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockManagePiiDetailsInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func TestManagePiiDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockManagePiiDetailsInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForManagePiiDetails(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostManagePii(t *testing.T) {
	assert := assert.New(t)
	client := &mockManagePiiDetailsInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManagePiiDetails(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/deputy/123?success=piiDetails"))
}

func TestErrorManagePiiMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockManagePiiDetailsInformation{}

	validationErrors := sirius.ValidationErrors{
		"piiReceived": {
			"isEmpty": "The pii received date is required and can't be empty",
		},
		"piiExpiry": {
			"isEmpty": "The pii expiry is required and can't be empty",
		},
		"piiAmount": {
			"isEmpty": "The pii amount is required and can't be empty",
		},
	}

	client.err = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/133", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForManagePiiDetails(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"piiReceived": {
				"isEmpty": "The pii received date is required and can't be empty",
			},
			"piiExpiry": {
				"isEmpty": "The pii expiry is required and can't be empty",
			},
			"piiAmount": {
				"isEmpty": "The pii amount is required and can't be empty",
			},
		},
	}

	assert.Equal(expectedValidationErrors, returnedError)
}

//write for expiry being before recived date
