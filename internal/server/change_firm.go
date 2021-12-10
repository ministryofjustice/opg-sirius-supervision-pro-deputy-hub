package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

func renderTemplateForChangeFirm(client ProDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:
			proDeputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)

			if err != nil {
				return err
			}

			firmDetails, err := client.GetFirms(ctx)

			if err != nil {
				return err
			}

			vars := proDeputyHubVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: proDeputyDetails,
				FirmDetails: firmDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			newFirm := r.PostFormValue("select-firm")

			AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")
			fmt.Println("existing firm id value")
			fmt.Println(AssignToExistingFirmStringIdValue)

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/deputy/%d/add-firm", deputyId))
			}

			AssignToFirmId := 0

			if AssignToExistingFirmStringIdValue != "" {
				AssignToFirmId, _ = strconv.Atoi(AssignToExistingFirmStringIdValue)
			}

			fmt.Println("AssignToFirmId")
			fmt.Println(AssignToFirmId)


			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, AssignToFirmId)

			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}

			//vars := proDeputyHubVars{
			//	Path:      r.URL.Path,
			//	XSRFToken: ctx.XSRFToken,
			//}

			return Redirect(fmt.Sprintf("/deputy/%d", deputyId))

			//return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
