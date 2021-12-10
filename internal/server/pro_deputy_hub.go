package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ProDeputyHubInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	GetFirms(sirius.Context) ([]sirius.Firm, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
}

type proDeputyHubVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	FirmDetails 	[]sirius.Firm
	Error            string
	Errors           sirius.ValidationErrors
	Success          bool
	SuccessMessage   string
	ExistingFirm bool
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

		hasSuccess := hasSuccessInUrl(r.URL.String(), "/deputy/"+strconv.Itoa(deputyId)+"/")

		vars := proDeputyHubVars{
			Path:             r.URL.Path,
			XSRFToken:        ctx.XSRFToken,
			ProDeputyDetails: proDeputyDetails,
			Success:          hasSuccess,
			SuccessMessage:   "Team details updated",
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}
