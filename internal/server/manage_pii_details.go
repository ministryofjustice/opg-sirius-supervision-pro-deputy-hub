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
}

func renderTemplateForManagePiiDetails(client ManagePiiDetailsInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			deputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
			if err != nil {
				return err
			}

			vars := proDeputyHubManagePiiVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: deputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			deputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)

			if err != nil {
				return err
			}

			addFirmDetailForm := sirius.PiiDetails{
				FirmId:       deputyDetails.Firm.FirmId,
				PiiReceived:  r.PostFormValue("pii-received"),
				PiiExpiry:    r.PostFormValue("pii-expiry"),
				PiiAmount:    r.PostFormValue("pii-amount"),
				PiiRequested: r.PostFormValue("pii-requested"),
			}

			err = client.EditPiiCertificate(ctx, addFirmDetailForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameEditPiiValidationErrorMessages(verr.Errors)
				vars := proDeputyHubManagePiiVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/deputy/%d?success=piiDetails", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameEditPiiValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if fieldName == "piiReceived" && errorType == "isEmpty" {
				err[errorType] = "The pii received date is required and can't be empty"
				errorCollection["piiReceived"] = err
			} else if fieldName == "piiExpiry" && errorType == "isEmpty" {
				err[errorType] = "The pii expiry is required and can't be empty"
				errorCollection["piiExpiry"] = err
			} else if fieldName == "piiAmount" && errorType == "isEmpty" {
				err[errorType] = "The pii amount is required and can't be empty"
				errorCollection["piiAmount"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
