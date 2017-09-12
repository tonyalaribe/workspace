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

// var (
// 	workPool *tunny.WorkPool
// )

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

func TriggerEvent(workspaceID, formID string, event database.TriggerEvent, data map[string]interface{}) {
	conf := config.Get()
	triggers, err := conf.Database.GetEventTriggers(workspaceID, formID, event)
	if err != nil {
		log.Println(err)
	}
	log.Println(triggers)
	for _, trigger := range triggers {
		PostToURL(trigger.URL, trigger.SecretToken, data)
	}

}

//
// func InitWorkerPool() {
// 	numCPUs := runtime.NumCPU()
// 	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.
// 	var err error
// 	workPool, err = tunny.CreatePoolGeneric(numCPUs).Open()
// 	if err != nil {
// 		log.Println(err)
// 	}
//
// 	defer workPool.Close()
// }
//
// func GetPool() *tunny.WorkPool {
// 	log.Printf("%#v", workPool)
// 	if workPool == nil {
// 		InitWorkerPool()
// 	}
// 	log.Printf("%#v", workPool)
// 	return workPool
// }
