package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ProDeputyChangeEcmInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	GetProTeamsMembers(sirius.Context, string, int) ([]sirius.TeamMember, error)
	ChangeECM(sirius.Context, sirius.ExecutiveCaseManagerOutgoing, sirius.ProDeputyDetails) error
}

type changeEcmVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	EcmTeamsDetails  []sirius.TeamMember
	Error            string
	Errors           sirius.ValidationErrors
	Success          bool
	SuccessMessage   string
}

func hasSuccessInUrl(url string, prefix string) bool {
	urlTrim := strings.TrimPrefix(url, prefix)
	return urlTrim == "?success=true"
}

func renderTemplateForChangeEcm(client ProDeputyChangeEcmInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
		if err != nil {
			return err
		}

		ecmTeamsDetails, err := client.GetProTeamsMembers(ctx, "PRO", proDeputyDetails.ExecutiveCaseManager.EcmId)
		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:
			var SuccessMessage string

			hasSuccess := hasSuccessInUrl(r.URL.String(), "/deputy/"+strconv.Itoa(deputyId))
			if hasSuccess {
				SuccessMessage = "new ecm is" + proDeputyDetails.ExecutiveCaseManager.EcmName
			}

			vars := changeEcmVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: proDeputyDetails,
				EcmTeamsDetails:  ecmTeamsDetails,
				Success:          hasSuccess,
				SuccessMessage:   SuccessMessage,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
			if err != nil {
				return err
			}

			vars := changeEcmVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: proDeputyDetails,
				EcmTeamsDetails:  ecmTeamsDetails,
			}

			EcmIdStringValue := r.PostFormValue("select-ecm")

			if EcmIdStringValue == "" {
				vars.Errors = sirius.ValidationErrors{
					"Change ECM": {"": "Select an executive case manager"},
				}
				EcmIdStringValue = "0"
			}

			EcmIdValue, err := strconv.Atoi(EcmIdStringValue)
			if err != nil {
				return err
			}

			changeECMForm := sirius.ExecutiveCaseManagerOutgoing{EcmId: EcmIdValue}

			err = client.ChangeECM(ctx, changeECMForm, proDeputyDetails)

			if len(vars.Errors) >= 1 {
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameManageDeputyContactDetailsValidationErrorMessages(verr.Errors)
				vars.Errors = verr.Errors

				return tmpl.ExecuteTemplate(w, "page", vars)
			}
			return Redirect(fmt.Sprintf("/deputy/%d?success=ecm", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
