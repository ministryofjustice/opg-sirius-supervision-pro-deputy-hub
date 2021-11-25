package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockProDeputyHubNotesInformation struct {
	count           int
	lastCtx         sirius.Context
	err             error
	addNote         error
	deputyData      sirius.ProDeputyDetails
	deputyNotesData sirius.DeputyNoteCollection
	userDetailsData sirius.UserDetails
}

func (m *mockProDeputyHubNotesInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockProDeputyHubNotesInformation) GetDeputyNotes(ctx sirius.Context, deputyId int) (sirius.DeputyNoteCollection, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyNotesData, m.err
}

func (m *mockProDeputyHubNotesInformation) AddNote(ctx sirius.Context, title, note string, deputyId, usedId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.addNote
}

func (m *mockProDeputyHubNotesInformation) GetUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userDetailsData, m.err
}

func TestGetNotes(t *testing.T) {
	assert := assert.New(t)

	client := &mockProDeputyHubNotesInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForProDeputyHubNotes(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal(2, client.count)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(proDeputyHubNotesVars{
		Path:           "/path",
		SuccessMessage: "Note added",
	}, template.lastVars)
}

func TestPostAddNote(t *testing.T) {
	assert := assert.New(t)
	client := &mockProDeputyHubNotesInformation{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForProDeputyHubNotes(client, nil)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(returnedError, Redirect("/deputy/123/notes?success=true"))
}

func TestErrorMessageWhenStringLengthTooLong(t *testing.T) {
	assert := assert.New(t)
	client := &mockProDeputyHubNotesInformation{}

	validationErrors := sirius.ValidationErrors{
		"name": {
			"stringLengthTooLong": "This team type is already in use",
		},
		"description": {
			"stringLengthTooLong": "This team type is already in use",
		},
	}
	client.addNote = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForProDeputyHubNotes(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"1-title": {
			"stringLengthTooLong": "The title must be 255 characters or fewer",
		},
		"2-note": {
			"stringLengthTooLong": "The note must be 1000 characters or fewer",
		},
	}

	assert.Equal(3, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		Path:   "/123",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}

func TestErrorMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockProDeputyHubNotesInformation{}

	validationErrors := sirius.ValidationErrors{
		"name": {
			"isEmpty": "This team type is already in use",
		},
		"description": {
			"isEmpty": "This team type is already in use",
		},
	}
	client.addNote = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForProDeputyHubNotes(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"1-title": {
			"isEmpty": "Enter a title for the note",
		},
		"2-note": {
			"isEmpty": "Enter a note",
		},
	}

	assert.Equal(3, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addNoteVars{
		Path:   "/123",
		Errors: expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
