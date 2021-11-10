package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ProDeputyHubInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
}

type proDeputyHubVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
}

func renderTemplateForProDeputyHub(client ProDeputyHubInformation, tmpl Template) Handler {
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

		vars := proDeputyHubVars{
			Path:             r.URL.Path,
			XSRFToken:        ctx.XSRFToken,
			ProDeputyDetails: proDeputyDetails,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
