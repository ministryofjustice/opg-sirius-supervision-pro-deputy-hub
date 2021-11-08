package server

import (
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"net/http"
)

type ProDeputyHubInformation interface {
}

type proDeputyHubVars struct {
	Path      string
	XSRFToken string
	Error     string
	Errors    sirius.ValidationErrors
}

func renderTemplateForProDeputyHub(client ProDeputyHubInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		ctx := getContext(r)

		vars := proDeputyHubVars{
			Path:      r.URL.Path,
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			return tmpl.ExecuteTemplate(w, "page", vars)
		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}