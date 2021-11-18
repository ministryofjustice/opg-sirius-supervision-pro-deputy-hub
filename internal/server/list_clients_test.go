package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeputyHubClientInformation struct {
	count            int
	lastCtx          sirius.Context
	err              error
	deputyData       sirius.ProDeputyDetails
	deputyClientData sirius.DeputyClientDetails
	ariaSorting      sirius.AriaSorting
}

func (m *mockDeputyHubClientInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockDeputyHubClientInformation) GetDeputyClients(ctx sirius.Context, deputyId int, columnBeingSorted string, sortOrder string) (sirius.DeputyClientDetails, sirius.AriaSorting, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyClientData, m.ariaSorting, m.err
}

func TestNavigateToClientTab(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeputyHubClientInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForClientTab(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestParseUrlReturnsColumnAndSortOrder(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/professional/deputy/78/clients?sort=crec:desc"
	expectedResponseColumnBeingSorted, sortOrder := "sort=crec", "desc"
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}

func TestParseUrlReturnsEmptyStrings(t *testing.T) {
	urlPassedin := "http://localhost:8888/supervision/deputies/professional/deputy/78/clients"
	expectedResponseColumnBeingSorted, sortOrder := "", ""
	resultColumnBeingSorted, resultSortOrder := parseUrl(urlPassedin)

	assert.Equal(t, expectedResponseColumnBeingSorted, resultColumnBeingSorted)
	assert.Equal(t, resultSortOrder, sortOrder)

}
