package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ManagePiiDetailsInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	EditPiiCertificate(sirius.Context, sirius.PiiDetails) error
}

type proDeputyHubManagePiiVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	DeputyNotes      sirius.DeputyNoteCollection
	Error            string
	Errors           sirius.ValidationErrors
	ErrorMessage     string
	Success          bool
	SuccessMessage   string
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		firmId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			deputyDetails, err := client.GetProDeputyDetails(ctx, firmId)
			if err != nil {
				return err
			}

			//hasSuccess := hasSuccessInUrl(r.URL.String(), "/deputy/"+strconv.Itoa(deputyId)+"/notes")

			vars := proDeputyHubManagePiiVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: deputyDetails,
				//Success:        hasSuccess,
				//SuccessMessage: "Note added",
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			addFirmDetailForm := sirius.PiiDetails{
				FirmId:       firmId,
				PiiReceived:  r.PostFormValue("pii-received"),
				PiiExpiry:    r.PostFormValue("pii-expiry"),
				PiiAmount:    r.PostFormValue("pii-amount"),
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			_, err := client.GetProDeputyDetails(ctx, deputyId)
			if err != nil {
				return err
			}

			err = client.EditPiiCertificate(ctx, addFirmDetailForm)

			// if verr, ok := err.(sirius.ValidationError); ok {

			// 	verr.Errors = renameValidationErrorMessages(verr.Errors)

			// 	vars = addNoteVars{
			// 		Path:             r.URL.Path,
			// 		XSRFToken:        ctx.XSRFToken,
			// 		Title:            title,
			// 		Note:             note,
			// 		Errors:           verr.Errors,
			// 		ProDeputyDetails: deputyDetails,
			// 	}

			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return tmpl.ExecuteTemplate(w, "page", vars)
			// } else
			if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/deputy/%d/notes?success=true", firmId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
