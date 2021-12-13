package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockProDeputyChangeFirmInformation struct {
	count      int
	lastCtx    sirius.Context
	err        error
	deputyData sirius.ProDeputyDetails
	firms []sirius.Firm
}

func (m *mockProDeputyChangeFirmInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockProDeputyChangeFirmInformation) GetFirms(ctx sirius.Context) ([]sirius.Firm, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.firms, m.err
}

func (m *mockProDeputyChangeFirmInformation) AssignDeputyToFirm(ctx sirius.Context, deputyId int, firmId int) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func TestChangeFirm(t *testing.T) {
	assert := assert.New(t)

	client := &mockProDeputyChangeFirmInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForProDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}


func TestCheckUrlForFirm(t *testing.T) {
	assert.Equal(t, false, checkUrlForFirm("http://localhost:8888/supervision/deputies/professional/deputy/76/change-firm"))
	assert.Equal(t, true, checkUrlForFirm("http://localhost:8888/supervision/deputies/professional/deputy/76/change-firm?existing-firm=true"))
}
