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

			vars := proDeputyHubVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: proDeputyDetails,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			newFirm := r.PostFormValue("select-firm")

			if newFirm == "new-firm" {
				return Redirect(fmt.Sprintf("/deputy/%d/add-firm", deputyId))
			}

			vars := proDeputyHubVars{
				Path:      r.URL.Path,
				XSRFToken: ctx.XSRFToken,
			}

			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
