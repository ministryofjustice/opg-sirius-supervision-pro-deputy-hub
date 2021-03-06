package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ProDeputyChangeFirmInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	GetFirms(sirius.Context) ([]sirius.Firm, error)
	AssignDeputyToFirm(sirius.Context, int, int) error
}

type changeFirmVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	FirmDetails      []sirius.Firm
	Error            string
	Errors           sirius.ValidationErrors
	Success          bool
	SuccessMessage   string
}

func renderTemplateForChangeFirm(client ProDeputyChangeFirmInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)

		if err != nil {
			return err
		}

		firmDetails, err := client.GetFirms(ctx)

		if err != nil {
			return err
		}

		switch r.Method {
		case http.MethodGet:

			vars := changeFirmVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: proDeputyDetails,
				FirmDetails:      firmDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars changeFirmVars
			newFirm := r.PostFormValue("select-firm")
			AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/deputy/%d/add-firm", deputyId))
			}

			AssignToFirmId := 0
			if AssignToExistingFirmStringIdValue != "" {
				AssignToFirmId, err = strconv.Atoi(AssignToExistingFirmStringIdValue)
				if err != nil {
					return err
				}
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, AssignToFirmId)

			if verr, ok := assignDeputyToFirmErr.(sirius.ValidationError); ok {

				verr.Errors = renameChangeFirmValidationErrorMessages(verr.Errors)

				vars = changeFirmVars{
					Path:      r.URL.Path,
					XSRFToken: ctx.XSRFToken,
					Errors:    verr.Errors,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}
			return Redirect(fmt.Sprintf("/deputy/%d?success=firm", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func checkUrlForFirm(url string) bool {
	splitStringByQuestion := strings.Split(url, "?")
	if len(splitStringByQuestion) > 1 {
		splitString := strings.Split(splitStringByQuestion[1], "=")
		return splitString[0] == "existing-firm"
	}
	return false
}

func renameChangeFirmValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)

			if fieldName == "firmId" && errorType == "notGreaterThanInclusive" {
				err[errorType] = "Enter a firm name or number"
				errorCollection["existing-firm"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
