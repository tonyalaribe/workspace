package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/actions"
	"gitlab.com/middlefront/workspace/database"
)

//NewFormSubmissionHandler is triggered when a new submission is made. This handler is saving the actual variable submission info
func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submission := database.SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		log.Println(err)
	}

	updatedSubmission, err := actions.NewFormSubmission(workspaceID, formID, submission)
	if err != nil {
		log.Println(err)
	}

	response := make(map[string]interface{})
	response["message"] = "Upload Success"
	response["submission"] = updatedSubmission
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

//UodateSubmissionHandler update submission
func UpdateSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	submissionID := httprouterParams.ByName("submissionID")

	newSubmission := database.SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&newSubmission)
	if err != nil {
		log.Println(err)
	}

	err = actions.UpdateSubmission(workspaceID, formID, submissionID, newSubmission)
	if err != nil {
		log.Println(err)
	}

	response := map[string]string{}
	response["message"] = "Upload Success"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

//GetSubmissionsHandler returns the submission based on formID
func GetSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submissions, err := actions.GetFormSubmissions(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(submissions)
	if err != nil {
		log.Println(err)
	}

}

//GetSubmissionInfoHandler returns information about given  submission
func GetSubmissionInfoHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	submissionID := httprouterParams.ByName("submissionID")

	submissionData, err := actions.GetFormSubmissionDetails(workspaceID, formID, submissionID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}
}

//DeleteSubmissionHandler  deletes a submission given submissionID and formID
func DeleteSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	submissionID := httprouterParams.ByName("submissionID")

	submissionData, err := actions.DeleteFormSubmission(workspaceID, formID, submissionID)

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}
}

//GetSubmissionChangelogHandler returns changelogs associated with submission
func GetSubmissionChangelogHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	submissionID := httprouterParams.ByName("submissionID")

	submissionData, err := actions.GetSubmissionChangelog(workspaceID, formID, submissionID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}

}
