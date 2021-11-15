package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ProDeputyHubTimelineInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	GetProDeputyTimeline(sirius.Context, int) (sirius.ProDeputyEventCollection, error)
}

type deputyHubEventVars struct {
	Path              string
	XSRFToken         string
	ProDeputyDetails  sirius.ProDeputyDetails
	ProDeputyTimeline sirius.ProDeputyEventCollection
	Error             string
	Errors            sirius.ValidationErrors
}

func renderTemplateForProDeputyHubTimeline(client ProDeputyHubTimelineInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
		if err != nil {
			return err
		}

		proDeputyTimeline, err := client.GetProDeputyTimeline(ctx, deputyId)
		if err != nil {
			return err
		}

		vars := deputyHubEventVars{
			Path:              r.URL.Path,
			XSRFToken:         ctx.XSRFToken,
			ProDeputyDetails:  proDeputyDetails,
			ProDeputyTimeline: proDeputyTimeline,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
