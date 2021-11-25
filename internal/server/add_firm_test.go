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

type mockFirmInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	addFirm          int
	deputToFirmyData error
}

func (m *mockFirmInformation) AddFirmDetails(ctx sirius.Context, deputyId sirius.FirmDetails) (int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.addFirm, m.err
}

func (m *mockFirmInformation) AssignDeputyToFirm(ctx sirius.Context, deputyId int, firmId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestGetFirm(t *testing.T) {
	assert := assert.New(t)

	client := &mockFirmInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForAddFirm(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(0, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
}

func TestPostAddFirm(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForAddFirm(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/deputy/123"))
}

func TestErrorEditDeputyMessageWhenStringLengthTooLong(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	validationErrors := sirius.ValidationErrors{
		"firmName": {
			"stringLengthTooLong": "What sirius gives us",
		}, "phoneNumber": {
			"stringLengthTooLong": "What sirius gives us",
		}, "email": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine1": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine2": {
			"stringLengthTooLong": "What sirius gives us",
		}, "addressLine3": {
			"stringLengthTooLong": "What sirius gives us",
		}, "town": {
			"stringLengthTooLong": "What sirius gives us",
		}, "county": {
			"stringLengthTooLong": "What sirius gives us",
		}, "postcode": {
			"stringLengthTooLong": "What sirius gives us",
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
		returnedError = renderTemplateForAddFirm(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"firmName": {
			"stringLengthTooLong": "The firm name must be 255 characters or fewer",
		}, "phoneNumber": {
			"stringLengthTooLong": "The telephone number must be 255 characters or fewer",
		}, "email": {
			"stringLengthTooLong": "The email must be 255 characters or fewer",
		}, "addressLine1": {
			"stringLengthTooLong": "The building or street must be 255 characters or fewer",
		}, "addressLine2": {
			"stringLengthTooLong": "Address line 2 must be 255 characters or fewer",
		}, "addressLine3": {
			"stringLengthTooLong": "Address line 3 must be 255 characters or fewer",
		}, "town": {
			"stringLengthTooLong": "The town or city must be 255 characters or fewer",
		}, "county": {
			"stringLengthTooLong": "The county must be 255 characters or fewer",
		}, "postcode": {
			"stringLengthTooLong": "The postcode must be 255 characters or fewer",
		},
	}

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addFirmVars{
		Path:   "/133",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorAddFirmMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockFirmInformation{}

	validationErrors := sirius.ValidationErrors{
		"firmName": {
			"isEmpty": "What sirius gives us",
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
		returnedError = renderTemplateForAddFirm(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"firmName": {
			"isEmpty": "The firm name is required and can't be empty",
		},
	}

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addFirmVars{
		Path:   "/133",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
