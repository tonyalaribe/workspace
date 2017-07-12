package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submission := database.SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		log.Println(err)
	}

	conf := config.Get()
	//Get the form metadata
	var formInfoByte []byte
	conf.DB.View(func(tx *bolt.Tx) error {
		formMetaBucket := tx.Bucket([]byte(conf.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(conf.FormsMetadata))
		formInfoByte = formMetaBucket.Get([]byte(formID))
		return nil
	})

	formInfo, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		log.Println(err)
	}

	schema := formInfo.Path("jsonschema")
	for k, v := range submission.FormData {
		schemaObject := schema.Path("properties").Search(k)
		switch schemaObject.Path("type").Data().(string) {
		case "string":
			itemFormat := ""
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				pathToItem, err := conf.FileManager.Save(submission.SubmissionName, workspaceID, v.(string))
				if err != nil {
					log.Println(err)
				}
				submission.FormData[k] = pathToItem
				break
			default:
				submission.FormData[k] = v.(string)
			}

		case "array":
			switch schemaObject.Path("items.type").Data().(string) {
			case "string":
				switch schemaObject.Path("items.format").Data().(string) {
				case "data-url":
					items := []string{}
					for _, item := range v.([]interface{}) {
						pathToItem, err := conf.FileManager.Save(submission.SubmissionName, workspaceID, item.(string))
						if err != nil {
							log.Println(err)
						}
						items = append(items, pathToItem)
					}
					submission.FormData[k] = items
				}
			}
		case "integer":
			//Using type float64 due to compiler complaints when handling integer types
			submission.FormData[k] = submission.FormData[k].(float64)
		default:
			submission.FormData[k] = submission.FormData[k].(string)
		}
	}

	conf.DB.Update(func(tx *bolt.Tx) error {
		formMetaBucket := tx.Bucket([]byte(conf.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(conf.FormsMetadata))

		formInfoByte = formMetaBucket.Get([]byte(formID))

		return nil
	})

	err = conf.Database.NewFormSubmission(workspaceID, formID, submission)
	if err != nil {
		log.Println(err)
	}
	byt, _ := json.MarshalIndent(submission, "", "\t")
	log.Println(string(byt))

	response := map[string]string{}
	response["message"] = "Upload Success"

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func UpdateSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}

	newSubmission := database.SubmissionData{}
	err = json.NewDecoder(r.Body).Decode(&newSubmission)
	if err != nil {
		log.Println(err)
	}

	//Get the previously updated data
	conf := config.Get()
	oldSubmission, err := conf.Database.GetFormSubmissionDetails(workspaceID, formID, submissionID)
	if err != nil {
		log.Println(err)
	}

	oldSubmission.Status = newSubmission.Status
	oldSubmission.LastModified = newSubmission.LastModified
	oldSubmission.FormData = newSubmission.FormData

	formInfoByte, err := conf.Database.GetFormJSONBySlug(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}
	formMetaData, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		log.Println(err)
	}

	schema := formMetaData.Path("jsonschema")
	for k, v := range newSubmission.FormData {
		schemaObject := schema.Path("properties").Search(k)
		switch schemaObject.Path("type").Data().(string) {
		case "string":
			itemFormat := ""
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				//file formatting
				pathToItem, err := conf.FileManager.Save(newSubmission.SubmissionName, workspaceID, v.(string))
				if err != nil {
					log.Println(err)
				}
				oldSubmission.FormData[k] = pathToItem
				break
			default:
				oldSubmission.FormData[k] = v.(string)
			}

		case "array":
			switch schemaObject.Path("items.type").Data().(string) {
			case "string":
				switch schemaObject.Path("items.format").Data().(string) {
				case "data-url":
					items := []string{}
					for _, item := range v.([]interface{}) {
						pathToItem, err := conf.FileManager.Save(newSubmission.SubmissionName, workspaceID, item.(string))
						if err != nil {
							log.Println(err)
						}
						items = append(items, pathToItem)
					}
					oldSubmission.FormData[k] = items
				}
			}
		case "integer":
			//Using type float64 due to compiler complaints when handling integer types
			oldSubmission.FormData[k] = newSubmission.FormData[k].(float64)
		default:
			oldSubmission.FormData[k] = newSubmission.FormData[k].(string)
		}
	}

	conf.Database.UpdateFormSubmission(workspaceID, formID, submissionID, oldSubmission)
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

func GetSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	conf := config.Get()
	submissions, err := conf.Database.GetFormSubmissions(workspaceID, formID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(submissions)
	if err != nil {
		log.Println(err)
	}

}

func GetSubmissionInfoHandler(w http.ResponseWriter, r *http.Request) {
	// username := r.Context().Value("username").(string)
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}
	conf := config.Get()
	submissionData, err := conf.Database.GetFormSubmissionDetails(workspaceID, formID, submissionID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}

}
