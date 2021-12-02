package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	ErrorHandlerClient
	ProDeputyHubInformation
	ProDeputyHubTimelineInformation
	ProDeputyHubClientInformation
	ProDeputyHubNotesInformation
	FirmInformation
	DeputyContactDetailsInformation
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, firmHubURL, webDir string) http.Handler {
	wrap := errorHandler(logger, client, templates["error.gotmpl"], prefix, siriusPublicURL)

	router := mux.NewRouter()
	router.Handle("/deputy/{id}",
		wrap(
			renderTemplateForProDeputyHub(client, templates["pro-dashboard.gotmpl"])))

	router.Handle("/deputy/{id}/clients",
		wrap(
			renderTemplateForClientTab(client, templates["clients.gotmpl"])))

	router.Handle("/deputy/{id}/timeline",
		wrap(
			renderTemplateForProDeputyHubTimeline(client, templates["timeline.gotmpl"])))

	router.Handle("/deputy/{id}/notes",
		wrap(
			renderTemplateForProDeputyHubNotes(client, templates["notes.gotmpl"])))

	router.Handle("/deputy/{id}/notes/add-note",
		wrap(
			renderTemplateForProDeputyHubNotes(client, templates["add-notes.gotmpl"])))
	router.Handle("/deputy/{id}/change-firm",
		wrap(
			renderTemplateForChangeFirm(client, templates["change-firm.gotmpl"])))

	router.Handle("/deputy/{id}/add-firm",
		wrap(
			renderTemplateForAddFirm(client, templates["add-firm.gotmpl"])))

	router.Handle("/deputy/{id}/manage-deputy-details",
		wrap(
			renderTemplateForManageDeputyDetails(client, templates["manage-deputy-details.gotmpl"])))

	router.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(webDir + "/static"))
	router.PathPrefix("/assets/").Handler(static)
	router.PathPrefix("/javascript/").Handler(static)
	router.PathPrefix("/stylesheets/").Handler(static)

	router.NotFoundHandler = notFoundHandler(templates["error.gotmpl"], siriusPublicURL)

	return http.StripPrefix(prefix, router)
}

type Redirect string

func (e Redirect) Error() string {
	return "redirect to " + string(e)
}

func (e Redirect) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	Firstname string
	Surname   string
	SiriusURL string
	Path      string
	Code      int
	Error     string
	Errors	  bool
}

type ErrorHandlerClient interface {
	MyPermissions(sirius.Context) (sirius.PermissionSet, error)
}

func errorHandler(logger Logger, client ErrorHandlerClient, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			myPermissions, err := client.MyPermissions(getContext(r))

			if err == nil {
				err = next(myPermissions, w, r)
			}

			if err != nil {
				if redirect, ok := err.(Redirect); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}

func notFoundHandler(tmplError Template, siriusURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmplError.ExecuteTemplate(w, "page", errorVars{
			SiriusURL: siriusURL,
			Code:      http.StatusNotFound,
			Error: "Not Found",
		})
	}
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
