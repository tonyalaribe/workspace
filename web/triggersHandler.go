package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/actions"
	"gitlab.com/middlefront/workspace/database"
)

//TriggerJSON encodes information about triggers and what events are active
type TriggerJSON struct {
	ID          string
	URL         string
	Endpoint    string
	SecretToken string

	NewSubmission     bool
	UpdateSubmission  bool
	ApproveSubmission bool
	DeleteSubmission  bool
}

//UpdateTriggerHandler Update a trigger in case of changes to the trigger
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
		err = actions.UpdateTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	} else {
		trigger.EventType = database.NewSubmissionTriggerEvent
		err = actions.DeleteTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	}

	if triggerJSON.UpdateSubmission {
		trigger.EventType = database.UpdateSubmissionTriggerEvent
		err = actions.UpdateTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	} else {
		trigger.EventType = database.UpdateSubmissionTriggerEvent
		err = actions.DeleteTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	}

	if triggerJSON.ApproveSubmission {
		trigger.EventType = database.ApproveSubmissionTriggerEvent
		err = actions.UpdateTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	} else {
		trigger.EventType = database.ApproveSubmissionTriggerEvent
		err = actions.DeleteTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	}

	if triggerJSON.DeleteSubmission {
		trigger.EventType = database.DeleteSubmissionTriggerEvent
		err = actions.UpdateTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	} else {
		trigger.EventType = database.DeleteSubmissionTriggerEvent
		err = actions.DeleteTrigger(trigger)
		if err != nil {
			log.Println(err)
		}
	}

	message := map[string]string{}
	message["message"] = "Strong"
	json.NewEncoder(w).Encode(message)
}

//GetFormTriggersHandler  Get all triggers associated with  a form
func GetFormTriggersHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	triggers, err := actions.GetFormTriggers(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}

	ts := make(map[string]TriggerJSON)
	for _, trigger := range triggers {
		ts[trigger.URL] = TriggerJSON{
			ID:          trigger.ID,
			URL:         trigger.URL, //TODO: delete one of URL and Endpoints
			SecretToken: trigger.SecretToken,
		}
	}

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
		} else if trigger.EventType == database.DeleteSubmissionTriggerEvent {
			t := ts[trigger.URL]
			t.DeleteSubmission = true
			ts[trigger.URL] = t
		}
	}

	triggersJSON := []TriggerJSON{}
	for _, v := range ts {
		triggersJSON = append(triggersJSON, v)
	}

	json.NewEncoder(w).Encode(triggersJSON)
}

//DeleteTriggerHandler Deletes a  trigger (removes all trigger event types)
func DeleteTriggerHandler(w http.ResponseWriter, r *http.Request) {
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

	trigger.EventType = database.NewSubmissionTriggerEvent
	err = actions.DeleteTrigger(trigger)
	if err != nil {
		log.Println(err)
	}

	trigger.EventType = database.UpdateSubmissionTriggerEvent
	err = actions.DeleteTrigger(trigger)
	if err != nil {
		log.Println(err)
	}

	trigger.EventType = database.ApproveSubmissionTriggerEvent
	err = actions.DeleteTrigger(trigger)
	if err != nil {
		log.Println(err)
	}

	trigger.EventType = database.DeleteSubmissionTriggerEvent
	err = actions.DeleteTrigger(trigger)
	if err != nil {
		log.Println(err)
	}

	message := map[string]string{}
	message["message"] = "Success"
	json.NewEncoder(w).Encode(message)
}

//TestTriggerHandler Send a test trigger to registered urls tied to trigger actions
func TestTriggerHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	triggerJSON := TriggerJSON{}
	err := json.NewDecoder(r.Body).Decode(&triggerJSON)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%#v", triggerJSON)
	demoSubmission := []byte(`
	{
		"formData":{
			"label1":"value1",
			"label2":"value2",
			"label3":"value3",
		},
		"created":1504252897414,
		"lastModified":1504252897414,
		"submissionName":"Submission Name",
		"status":"draft",
		"id":24,
		"submissionNotes":"lorem ipsum dolores "
	}
`)

	demoChangelog := []byte(`
	{
		"created":1504252897414,
		"lastModified":1504252897414,
		"submissionName":"Submission Name",
		"submissionNotes":"lorem ipsum dolores "
	}
`)

	if triggerJSON.NewSubmission {
		data := make(map[string]interface{})
		data["workspaceID"] = workspaceID
		data["formID"] = formID
		data["event"] = string(database.NewSubmissionTriggerEvent)
		data["submission"] = demoSubmission

		actions.TriggerEvent(workspaceID, formID, database.NewSubmissionTriggerEvent, data)

		if err != nil {
			log.Println(err)
		}
	}

	if triggerJSON.UpdateSubmission {

		data := make(map[string]interface{})
		data["workspaceID"] = workspaceID
		data["formID"] = formID
		data["event"] = string(database.UpdateSubmissionTriggerEvent)
		data["submission"] = demoSubmission
		data["changelog"] = demoChangelog

		actions.TriggerEvent(workspaceID, formID, database.UpdateSubmissionTriggerEvent, data)

		if err != nil {
			log.Println(err)
		}

	}

	if triggerJSON.ApproveSubmission {

		data := make(map[string]interface{})
		data["workspaceID"] = workspaceID
		data["formID"] = formID
		data["event"] = string(database.ApproveSubmissionTriggerEvent)
		data["submission"] = demoSubmission

		actions.TriggerEvent(workspaceID, formID, database.DeleteSubmissionTriggerEvent, data)

		if err != nil {
			log.Println(err)
		}
	}

	if triggerJSON.DeleteSubmission {

		data := make(map[string]interface{})
		data["workspaceID"] = workspaceID
		data["formID"] = formID
		data["event"] = string(database.DeleteSubmissionTriggerEvent)
		data["submission"] = demoSubmission

		actions.TriggerEvent(workspaceID, formID, database.DeleteSubmissionTriggerEvent, data)

		if err != nil {
			log.Println(err)
		}
	}

	message := map[string]string{}
	message["message"] = "Strong"
	json.NewEncoder(w).Encode(message)
}
