package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type ProDeputyHubNotesInformation interface {
	GetProDeputyDetails(sirius.Context, int) (sirius.ProDeputyDetails, error)
	GetDeputyNotes(sirius.Context, int) (sirius.DeputyNoteCollection, error)
	AddNote(ctx sirius.Context, title, note string, deputyId, userId int) error
	GetUserDetails(sirius.Context) (sirius.UserDetails, error)
}

type proDeputyHubNotesVars struct {
	Path             string
	XSRFToken        string
	ProDeputyDetails sirius.ProDeputyDetails
	DeputyNotes      sirius.DeputyNoteCollection
	Error            string
	Errors           sirius.ValidationErrors
	ErrorMessage     string
	Success          bool
	SuccessMessage   string
}

type addNoteVars struct {
	Path             string
	XSRFToken        string
	Title            string
	Note             string
	Success          bool
	Error            sirius.ValidationError
	Errors           sirius.ValidationErrors
	ErrorMessage     string
	ProDeputyDetails sirius.ProDeputyDetails
}

func hasSuccessInUrl(url string, prefix string) bool {
	urlTrim := strings.TrimPrefix(url, prefix)
	return urlTrim == "?success=true"
}

func renderTemplateForProDeputyHubNotes(client ProDeputyHubNotesInformation, tmpl Template) Handler {
	return func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error {

		ctx := getContext(r)
		routeVars := mux.Vars(r)
		deputyId, _ := strconv.Atoi(routeVars["id"])

		switch r.Method {
		case http.MethodGet:

			deputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
			if err != nil {
				return err
			}
			deputyNotes, err := client.GetDeputyNotes(ctx, deputyId)
			if err != nil {
				return err
			}

			hasSuccess := hasSuccessInUrl(r.URL.String(), "/deputy/"+strconv.Itoa(deputyId)+"/notes")

			vars := proDeputyHubNotesVars{
				Path:             r.URL.Path,
				XSRFToken:        ctx.XSRFToken,
				ProDeputyDetails: deputyDetails,
				DeputyNotes:      deputyNotes,
				Success:          hasSuccess,
				SuccessMessage:   "Note added",
			}

			return tmpl.ExecuteTemplate(w, "page", vars)

		case http.MethodPost:
			var vars addNoteVars
			var (
				title = r.PostFormValue("title")
				note  = r.PostFormValue("note")
			)

			userId, err := client.GetUserDetails(ctx)
			if err != nil {
				return err
			}

			deputyDetails, err := client.GetProDeputyDetails(ctx, deputyId)
			if err != nil {
				return err
			}

			err = client.AddNote(ctx, title, note, deputyId, userId.ID)

			if verr, ok := err.(sirius.ValidationError); ok {

				verr.Errors = renameValidationErrorMessages(verr.Errors)

				vars = addNoteVars{
					Path:             r.URL.Path,
					XSRFToken:        ctx.XSRFToken,
					Title:            title,
					Note:             note,
					Errors:           verr.Errors,
					ProDeputyDetails: deputyDetails,
				}

				w.WriteHeader(http.StatusBadRequest)
				return tmpl.ExecuteTemplate(w, "page", vars)
			} else if err != nil {
				return err
			}

			return Redirect(fmt.Sprintf("/deputy/%d/notes?success=true", deputyId))

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}

func renameValidationErrorMessages(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	errorCollection := sirius.ValidationErrors{}

	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if fieldName == "name" && errorType == "stringLengthTooLong" {
				err[errorType] = "The title must be 255 characters or fewer"
				errorCollection["1-title"] = err
			} else if fieldName == "name" && errorType == "isEmpty" {
				err[errorType] = "Enter a title for the note"
				errorCollection["1-title"] = err
			} else if fieldName == "description" && errorType == "stringLengthTooLong" {
				err[errorType] = "The note must be 1000 characters or fewer"
				errorCollection["2-note"] = err
			} else if fieldName == "description" && errorType == "isEmpty" {
				err[errorType] = "Enter a note"
				errorCollection["2-note"] = err
			} else {
				err[errorType] = errorMessage
				errorCollection[fieldName] = err
			}
		}
	}
	return errorCollection
}
