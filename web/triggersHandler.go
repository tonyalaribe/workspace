package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/actions"
	"gitlab.com/middlefront/workspace/database"
)

type TriggerJSON struct {
	ID          string
	URL         string
	Endpoint    string
	SecretToken string

	NewSubmission     bool
	UpdateSubmission  bool
	ApproveSubmission bool
}

func UpdateTriggerHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	triggerJSON := TriggerJSON{}
	err := json.NewDecoder(r.Body).Decode(&triggerJSON)
	if err != nil {
		log.Println(err)
	}

	trigger := database.Trigger{
		ID:          triggerJSON.ID,
		WorkspaceID: workspaceID,
		FormID:      formID,
		URL:         triggerJSON.URL,
		SecretToken: triggerJSON.SecretToken,
	}

	if triggerJSON.NewSubmission {
		trigger.EventType = database.NewSubmissionTriggerEvent
	} else if triggerJSON.UpdateSubmission {
		trigger.EventType = database.UpdateSubmissionTriggerEvent
	} else if triggerJSON.ApproveSubmission {
		trigger.EventType = database.ApproveSubmissionTriggerEvent
	} else {
		trigger.EventType = database.NewSubmissionTriggerEvent
	}
	log.Println("update trigger")
	log.Println(trigger)
	err = actions.UpdateTrigger(trigger)
	if err != nil {
		log.Println(err)
	}
}

func GetTriggersHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	url := r.URL.Query().Get("url")

	triggerJSON := TriggerJSON{}
	err := json.NewDecoder(r.Body).Decode(&triggerJSON)
	if err != nil {
		log.Println(err)
	}

	resp, err := actions.GetTriggers(workspaceID, formID, url)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(resp)
}
