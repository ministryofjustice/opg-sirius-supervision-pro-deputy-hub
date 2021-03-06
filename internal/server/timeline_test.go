package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockProDeputyHubTimelineInformation struct {
	count        int
	lastCtx      sirius.Context
	err          error
	deputyData   sirius.ProDeputyDetails
	timelineData sirius.ProDeputyEventCollection
}

func (m *mockProDeputyHubTimelineInformation) GetProDeputyDetails(ctx sirius.Context, deputyId int) (sirius.ProDeputyDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.deputyData, m.err
}

func (m *mockProDeputyHubTimelineInformation) GetProDeputyTimeline(ctx sirius.Context, deputyId int) (sirius.ProDeputyEventCollection, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.timelineData, m.err
}

func TestNavigateToTimeline(t *testing.T) {
	assert := assert.New(t)

	client := &mockProDeputyHubTimelineInformation{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := renderTemplateForProDeputyHub(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}
