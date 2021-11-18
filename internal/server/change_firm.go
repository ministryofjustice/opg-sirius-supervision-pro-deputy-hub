package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
)

type editProDeputyFirmVars struct {
	Path      string
	XSRFToken string
	Error     string
	Errors    sirius.ValidationErrors
	Success   bool
}

func renderTemplateForChangeFirm(client ProDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		return Redirect(fmt.Sprintf("/deputy/%d/add-firm", deputyId))
	}
}
