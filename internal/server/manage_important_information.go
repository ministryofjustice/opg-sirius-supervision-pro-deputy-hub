package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type ManageDeputyImportantInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	UpdateImportantInformation(sirius.Context, int, sirius.ImportantInformationDetails) error
	GetDeputyAnnualInvoiceBillingTypes(ctx sirius.Context) ([]sirius.DeputyAnnualBillingInvoiceTypes, error)
}

type manageDeputyImportantInformationVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	Error            string
	Errors           sirius.ValidationErrors
	DeputyId         int
	AnnualBillingInvoiceTypes []sirius.DeputyAnnualBillingInvoiceTypes
}

func renderTemplateForImportantInformation(client ManageDeputyImportantInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)

		if err != nil {
			return err
		}

		annualBillingInvoiceTypes, err := client.GetDeputyAnnualInvoiceBillingTypes(ctx)
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
				AnnualBillingInvoiceTypes: annualBillingInvoiceTypes,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)


		case http.MethodPost:
			var complaintsBool bool

			panelDeputyBool, err := strconv.ParseBool(r.PostFormValue("panel-deputy"))
			if err != nil{
				return err
			}

			importantInfoForm := sirius.ImportantInformationDetails{
				Complaints:  complaintsBool,
				PanelDeputy:  panelDeputyBool,
				AnnualBillingInvoice: r.PostFormValue("annual-billing"),
				OtherImportantInformation:     r.PostFormValue("other-information"),
			}

			fmt.Println("imp info form abi")
			fmt.Println(importantInfoForm.AnnualBillingInvoice)

			err = client.UpdateImportantInformation(ctx, deputyId, importantInfoForm)

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

			return Redirect(fmt.Sprintf("/deputy/%d?success=importantInformation", deputyId))
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
