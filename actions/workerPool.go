package actions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

//PostToURL sends event info to registered webhook
func PostToURL(url string, secretToken string, body interface{}) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("X-Workspace-Token", secretToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)

	} else {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Println(respBody)
	}

}

//TriggerEvent triggers relevant webhooks for given event
func TriggerEvent(workspaceID, formID string, event database.TriggerEvent, data map[string]interface{}) {
	conf := config.Get()
	triggers, err := conf.Database.GetEventTriggers(workspaceID, formID, event)
	if err != nil {
		log.Println(err)
	}
	for _, trigger := range triggers {
		PostToURL(trigger.URL, trigger.SecretToken, data)
	}

}
