package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type DeputyContactDetailsInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	UpdateDeputyContactDetails(sirius.Context, int, sirius.DeputyContactDetails) error
}

type manageDeputyDetailsVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
	Success          bool
	DeputyId         int
}

func renderTemplateForManageDeputyDetails(client DeputyContactDetailsInformation, tmpl Template) Handler {
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

			vars := manageDeputyDetailsVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				DeputyId:         deputyId,
				ProDeputyDetails: proDeputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			form := sirius.DeputyContactDetails{
				DeputyFirstName: r.PostFormValue("deputy-first-name"),
				DeputySurname:   r.PostFormValue("deputy-last-name"),
				AddressLine1:    r.PostFormValue("address-line-1"),
				AddressLine2:    r.PostFormValue("address-line-2"),
				AddressLine3:    r.PostFormValue("address-line-3"),
				Town:            r.PostFormValue("town"),
				County:          r.PostFormValue("county"),
				Postcode:        r.PostFormValue("postcode"),
				PhoneNumber:     r.PostFormValue("telephone"),
				Email:           r.PostFormValue("email"),
			}

			err := client.UpdateDeputyContactDetails(ctx, deputyId, form)

			if verr, ok := err.(sirius.ValidationError); ok {
				verr.Errors = renameManageDeputyDetailsValidationErrorMessages(verr.Errors)
				vars := manageDeputyDetailsVars{
					Path:             r.URL.Path,
					XSRFToken:        ctx.XSRFToken,
					DeputyId:         deputyId,
					ProDeputyDetails: proDeputyDetails,
					Errors:           verr.Errors,
				}
				return tmpl.ExecuteTemplate(w, "page", vars)
			}

			return Redirect(fmt.Sprintf("/deputy/%d?success=deputyDetails", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameManageDeputyDetailsValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)

			if fieldName == "organisationName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The deputy name must be 255 characters or fewer"
				errorCollection["organisationName"] = err
			} else if fieldName == "organisationName" && errorType == "isEmpty" {
				err[errorType] = "Enter a deputy name"
				errorCollection["organisationName"] = err
			} else if fieldName == "workPhoneNumber" && errorType == "stringLengthTooLong" {
				err[errorType] = "The telephone number must be 255 characters or fewer"
				errorCollection["workPhoneNumber"] = err
			} else if fieldName == "email" && errorType == "stringLengthTooLong" {
				err[errorType] = "The email number must be 255 characters or fewer"
				errorCollection["email"] = err
			} else if fieldName == "organisationTeamOrDepartmentName" && errorType == "stringLengthTooLong" {
				err[errorType] = "The team or department must be 255 characters or fewer"
				errorCollection["organisationTeamOrDepartmentName"] = err
			} else if fieldName == "addressLine1" && errorType == "stringLengthTooLong" {
				err[errorType] = "The building or street must be 255 characters or fewer"
				errorCollection["addressLine1"] = err
			} else if fieldName == "addressLine2" && errorType == "stringLengthTooLong" {
				err[errorType] = "Address line 2 must be 255 characters or fewer"
				errorCollection["addressLine2"] = err
			} else if fieldName == "addressLine3" && errorType == "stringLengthTooLong" {
				err[errorType] = "AddressLine 3 must be 255 characters or fewer"
				errorCollection["addressLine3"] = err
			} else if fieldName == "town" && errorType == "stringLengthTooLong" {
				err[errorType] = "The town or city must be 255 characters or fewer"
				errorCollection["town"] = err
			} else if fieldName == "county" && errorType == "stringLengthTooLong" {
				err[errorType] = "The county must be 255 characters or fewer"
				errorCollection["county"] = err
			} else if fieldName == "postcode" && errorType == "stringLengthTooLong" {
				err[errorType] = "The postcode must be 255 characters or fewer"
				errorCollection["postcode"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
