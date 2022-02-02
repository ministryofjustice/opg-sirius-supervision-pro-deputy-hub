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

type mockManageDeputyImportantInformation struct {
	count                     int
	lastCtx                   sirius.Context
	err                       error
	deputyData                sirius.ProDeputyDetails
	updateErr                 error
	annualBillingInvoiceTypes []sirius.DeputyAnnualBillingInvoiceTypes
	complaintTypes            []sirius.DeputyComplaintTypes
}

func (m *mockManageDeputyImportantInformation) GetProDeputyDetails(ctx sirius.Context, _ int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockManageDeputyImportantInformation) GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.annualBillingInvoiceTypes, m.err
}

func (m *mockManageDeputyImportantInformation) GetDeputyComplaintTypes(ctx sirius.Context) ([]sirius.DeputyComplaintTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.complaintTypes, m.err
}

func (m *mockManageDeputyImportantInformation) UpdateImportantInformation(ctx sirius.Context, _ int, _ sirius.ImportantInformationDetails) error {
	m.count += 1
	m.lastCtx = ctx

	return m.updateErr
}

func TestGetManageImportantInformation(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyImportantInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := renderTemplateForImportantInformation(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestPostManageImportantInformation(t *testing.T) {
	assert := assert.New(t)

	client := &mockManageDeputyImportantInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var redirect error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		redirect = renderTemplateForImportantInformation(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)
	assert.Equal(Redirect("/deputy/123?success=importantInformation"), redirect)
}

func TestErrorManageImportantInformationMessageWhenIsEmpty(t *testing.T) {
	assert := assert.New(t)
	client := &mockManageDeputyImportantInformation{}

	validationErrors := sirius.ValidationErrors{
		"otherImportantInformation": {
			"stringLengthTooLong": "What sirius gives us",
		},
	}

	client.updateErr = sirius.ValidationError{
		Errors: validationErrors,
	}

	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/123", strings.NewReader(""))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var returnedError error

	testHandler := mux.NewRouter()
	testHandler.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		returnedError = renderTemplateForImportantInformation(client, template)(sirius.PermissionSet{}, w, r)
	})

	testHandler.ServeHTTP(w, r)

	expectedValidationErrors := sirius.ValidationErrors{
		"otherImportantInformation": {
			"stringLengthTooLong": "The other important information must be 1000 characters or fewer",
		},
	}

	assert.Equal(4, client.count)

	assert.Equal(manageDeputyImportantInformationVars{
		Path:     "/123",
		DeputyId: 123,
		Errors:   expectedValidationErrors,
	}, template.lastVars)

	assert.Nil(returnedError)
}
