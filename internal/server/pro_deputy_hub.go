package server

import (
	"net/http"
	"strconv"
	"strings"

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
	Success          bool
	SuccessMessage   string
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

		hasSuccess, successMessage := createSuccessAndSuccessMessageForVars(r.URL.String(), proDeputyDetails.Firm.FirmName)

		vars := proDeputyHubVars{
			Path:             r.URL.Path,
			XSRFToken:        ctx.XSRFToken,
			ProDeputyDetails: proDeputyDetails,
			Success:          hasSuccess,
			SuccessMessage:   successMessage,
		}

		return tmpl.ExecuteTemplate(w, "page", vars)
	}
}

func createSuccessAndSuccessMessageForVars(url, firmName string) (bool, string) {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")

		if splitString[1] == "firm" {
			return true, "Firm changed to " + firmName
		}
	}
	return false, ""
}