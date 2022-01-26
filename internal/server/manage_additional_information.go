package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type DeputyImportantInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	UpdateAdditionalInformation(sirius.Context, int, sirius.AdditionalInformationDetails) error
}

type manageDeputyImportantInformationVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
	DeputyId         int
}

func renderTemplateForAdditionalInformation(client DeputyImportantInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)

		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			vars := manageDeputyImportantInformationVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				DeputyId:         deputyId,
				ProDeputyDetails: proDeputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			additionalInfoForm := sirius.AdditionalInformationDetails{
				//DeputySubType:    proDeputyDetails.DeputySubType.SubType,
				//DeputyFirstName:  r.PostFormValue("deputy-first-name"),
				//DeputySurname:    r.PostFormValue("deputy-last-name"),
				//OrganisationName: r.PostFormValue("organisation-name"),
				//AddressLine1:     r.PostFormValue("address-line-1"),
				//AddressLine2:     r.PostFormValue("address-line-2"),
				//AddressLine3:     r.PostFormValue("address-line-3"),
				//Town:             r.PostFormValue("town"),
				//County:           r.PostFormValue("county"),
				//Postcode:         r.PostFormValue("postcode"),
				//PhoneNumber:      r.PostFormValue("telephone"),
				//Email:            r.PostFormValue("email"),
			}

			err := client.UpdateAdditionalInformation(ctx, deputyId, additionalInfoForm)

			if verr, ok := err.(sirius.ValidationError); ok {
				//verr.Errors = renameUpdateAdditionalInformationValidationErrorMessages(verr.Errors)
				vars := manageDeputyImportantInformationVars{
					Path:             r.URL.Path,
					XSRFToken:        ctx.XSRFToken,
					DeputyId:         deputyId,
					ProDeputyDetails: proDeputyDetails,
					Errors:           verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/deputy/%d?success=additionalInformation", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

//func renameUpdateAdditionalInformationValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
//	errorCollection := sirius.ValidationErrors{}
//
//	for fieldName, value := range siriusError {
//		for errorType, errorMessage := range value {
//			err := make(map[string]string)
//
//			if fieldName == "firstname" && errorType == "stringLengthTooLong" {
//				err[errorType] = "The deputy first name must be 255 characters or fewer"
//				errorCollection["firstname"] = err
//			} else {
//				err[errorType] = errorMessage
//				errorCollection[fieldName] = err
//			}
//		}
//	}
//	return errorCollection
//}
