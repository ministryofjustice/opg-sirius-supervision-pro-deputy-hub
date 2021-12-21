package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockProDeputyHubInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.ProDeputyDetails
}

func (m *mockProDeputyHubInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func TestNavigateToProDeputyHub(t *testing.T) {
	assert := assert.New(t)

	client := &mockProDeputyHubInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForProDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestCreateSuccessAndSuccessMessageForVars(t *testing.T) {
	tests := []struct {
		url                    string
		firmname               string
		expectedBool           bool
		expectedSuccessMessage string
	}{
		{url: "localhost:8888/deputies/1?success=firm", firmname: "lovely firm", expectedBool: true, expectedSuccessMessage: "Firm changed to lovely firm"},
		{url: "localhost:8888/deputies/1", firmname: "lovely firm", expectedBool: false, expectedSuccessMessage: ""},
		{url: "localhost:8888/deputies/1?success=newFirm", firmname: "interesting firm", expectedBool: true, expectedSuccessMessage: "Firm added"},
		{url: "localhost:8888/deputies/1?success=deputyDetails", firmname: "another firm", expectedBool: true, expectedSuccessMessage: "Deputy details updated"},
		{url: "localhost:8888/deputies/1?success=piiDetails", firmname: "firm number 4", expectedBool: true, expectedSuccessMessage: "Pii details updated"},
	}

	for _, tc := range tests {
		returnedBool, returnedString := createSuccessAndSuccessMessageForVars(tc.url, tc.firmname)
		assert.Equal(t, tc.expectedBool, returnedBool)
		assert.Equal(t, tc.expectedSuccessMessage, returnedString)
	}
}
