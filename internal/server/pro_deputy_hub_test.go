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

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnEcmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=ecm", "firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Ecm changed to Jon Snow")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnDeputyDetailsSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=deputyDetails", "firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Deputy details updated")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageOnChangeFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=firm", "firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm changed to firm Name")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsMessageAddFirmSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=newFirm", "firm Name", "Jon Snow")
	assert.Equal(t, true, Success)
	assert.Equal(t, SuccessMessage, "Firm added")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilForAnyOtherText(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/?success=otherMessage", "firm Name", "Jon Snow")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}

func TestCreateSuccessAndSuccessMessageForVarsReturnsNilIfNoSuccess(t *testing.T) {
	Success, SuccessMessage := createSuccessAndSuccessMessageForVars("/deputy/76/", "firm Name", "Jon Snow")
	assert.Equal(t, false, Success)
	assert.Equal(t, SuccessMessage, "")
}
