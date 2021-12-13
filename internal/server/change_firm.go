package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

func renderTemplateForChangeFirm(client ProDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])
		firm := checkUrlForFirm(r.URL.String())

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
				ExistingFirm: firm,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			newFirm := r.PostFormValue("select-firm")
			AssignToExistingFirmStringIdValue := r.PostFormValue("select-existing-firm")

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/deputy/%d/add-firm", deputyId))
			}

			AssignToFirmId := 0

			if AssignToExistingFirmStringIdValue != "" {
				AssignToFirmId, _ = strconv.Atoi(AssignToExistingFirmStringIdValue)
			}

			assignDeputyToFirmErr := client.AssignDeputyToFirm(ctx, deputyId, AssignToFirmId)

			if assignDeputyToFirmErr != nil {
				return assignDeputyToFirmErr
			}
			return Redirect(fmt.Sprintf("/deputy/%d", deputyId))
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

	func checkUrlForFirm(url string) bool {
	splitStringByQuestion := strings.Split(url, "?")
		if len(splitStringByQuestion) > 1 {
			splitString := strings.Split(splitStringByQuestion[1], "=")
			if splitString[1] == "existing-firm" {
				return true
			}
			return false
		}
		return false
	}

