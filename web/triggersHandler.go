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

	err = actions.UpdateTrigger(trigger)
	if err != nil {
		log.Println(err)
	}
}

func GetFormTriggersHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	// url := r.URL.Query().Get("url")

	// triggerJSON := TriggerJSON{}

	triggers, err := actions.GetFormTriggers(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}

	log.Println(triggers)
	ts := make(map[string]TriggerJSON)
	for _, trigger := range triggers {
		ts[trigger.URL] = TriggerJSON{
			ID:          trigger.ID,
			URL:         trigger.URL, //TODO: delete one of URL and Endpoints
			SecretToken: trigger.SecretToken,
		}
	}
	log.Println(ts)
	for _, trigger := range triggers {
		if trigger.EventType == database.NewSubmissionTriggerEvent {
			t := ts[trigger.URL]
			t.NewSubmission = true
			ts[trigger.URL] = t
		} else if trigger.EventType == database.UpdateSubmissionTriggerEvent {
			t := ts[trigger.URL]
			t.UpdateSubmission = true
			ts[trigger.URL] = t
		} else if trigger.EventType == database.ApproveSubmissionTriggerEvent {
			t := ts[trigger.URL]
			t.ApproveSubmission = true
			ts[trigger.URL] = t
		}
	}

	triggersJSON := []TriggerJSON{}
	for _, v := range ts {
		triggersJSON = append(triggersJSON, v)
	}

	json.NewEncoder(w).Encode(triggersJSON)
}