package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/config"
)

type SubmissionData struct {
	FormData       map[string]interface{} `json:"formData"`
	Created        int                    `json:"created"`
	LastModified   int                    `json:"lastModified"`
	SubmissionName string                 `json:"submissionName"`
	Status         string                 `json:"status"`
	ID             int                    `json:"id"`
}

type Files struct {
	Status     string `json:"status"`
	File       string `json:"file"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Path       string `json:"path"`
	CreatedBy  string `json:"createdBy"`
	UploadDate string `json:"uploadDate"`
}

func Base64ToFileSystem(b64 string, location string) string {
	if len(strings.Split(b64, "base64,")) < 2 {
		return b64
	}
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
		//Cant decode string, so its probably an already processed url.
		return b64
	}

	meta := strings.Split(b64, "base64,")[0]
	filename := strings.Replace(strings.Split(meta, "name=")[1], ";", "", -1)

	fullPath := filepath.Join(location, filename)
	err = ioutil.WriteFile(fullPath, byt, 0644)
	if err != nil {
		log.Println(err)
	}
	return fullPath
}

func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submission := SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", submission)

	conf := config.Get()

	//Get the form metadata
	var formInfoByte []byte
	conf.DB.View(func(tx *bolt.Tx) error {
		formMetaBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(config.FORMS_METADATA))

		formInfoByte = formMetaBucket.Get([]byte(formID))

		return nil
	})
	formInfo, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", formInfo)

	// schema := formInfo["jsonschema"].(map[string]interface{})
	schema := formInfo.Path("jsonschema")

	for k, v := range submission.FormData {
		schemaObject := schema.Path("properties").Search(k)

		switch schemaObject.Path("type").Data().(string) {
		case "string":

			itemFormat := ""
			// if schemaObject["format"] != nil {
			// 	itemFormat = schemaObject["format"].(string)
			// }
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				//file formatting
				pathToSubmission := filepath.Join(conf.RootDirectory, username, submission.SubmissionName)
				os.MkdirAll(pathToSubmission, os.ModePerm)
				fullPath := Base64ToFileSystem(v.(string), pathToSubmission)
				submission.FormData[k] = fullPath
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
						log.Println("data-uri processing")
						pathToSubmission := filepath.Join(conf.RootDirectory, username, submission.SubmissionName)
						os.MkdirAll(pathToSubmission, os.ModePerm)
						fullPath := Base64ToFileSystem(item.(string), pathToSubmission)
						items = append(items, fullPath)
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
		formMetaBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(config.FORMS_METADATA))

		formInfoByte = formMetaBucket.Get([]byte(formID))

		return nil
	})
	/*Save to boltdb*/
	err = conf.DB.Update(func(tx *bolt.Tx) error {

		formBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		nextID, err := formBucket.NextSequence()
		if err != nil {
			log.Println(err)
		}

		submission.ID = int(nextID)
		dataByte, err := json.Marshal(submission)
		if err != nil {
			log.Println(err)
		}

		err = formBucket.Put(itob(int(nextID)), dataByte)
		if err != nil {
			log.Println(err)
		}

		return nil
	})

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
	username := r.Context().Value("username").(string)

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}

	newSubmission := SubmissionData{}
	err = json.NewDecoder(r.Body).Decode(&newSubmission)
	if err != nil {
		log.Println(err)
	}

	//Get the previously updated data
	oldSubmission := SubmissionData{}
	conf := config.Get()
	err = conf.DB.View(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		err = json.Unmarshal(formBucket.Get(itob(submissionID)), &oldSubmission)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	oldSubmission.Status = newSubmission.Status
	oldSubmission.LastModified = newSubmission.LastModified
	oldSubmission.FormData = newSubmission.FormData

	//Get the form meta data. SO
	var formInfoByte []byte
	conf.DB.View(func(tx *bolt.Tx) error {
		formMetaBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(config.FORMS_METADATA))

		formInfoByte = formMetaBucket.Get([]byte(formID))
		return nil
	})
	formMetaData, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		log.Println(err)
	}

	schema := formMetaData.Path("jsonschema")
	for k, v := range newSubmission.FormData {

		schemaObject := schema.Path("properties").Search(k)
		log.Println(schemaObject)

		switch schemaObject.Path("type").Data().(string) {
		case "string":
			itemFormat := ""
			// if schemaObject["format"] != nil {
			// 	itemFormat = schemaObject["format"].(string)
			// }
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				//file formatting
				pathToSubmission := filepath.Join(conf.RootDirectory, username, newSubmission.SubmissionName)
				os.MkdirAll(pathToSubmission, os.ModePerm)
				fullPath := Base64ToFileSystem(v.(string), pathToSubmission)
				oldSubmission.FormData[k] = fullPath
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
						log.Println("data-uri processing")
						pathToSubmission := filepath.Join(conf.RootDirectory, username, newSubmission.SubmissionName)
						os.MkdirAll(pathToSubmission, os.ModePerm)
						fullPath := Base64ToFileSystem(item.(string), pathToSubmission)
						items = append(items, fullPath)
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

	dataByte, err := json.Marshal(oldSubmission)
	if err != nil {
		log.Println(err)
	}
	err = conf.DB.Update(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		err = formBucket.Put(itob(submissionID), dataByte)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

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
	// username := r.Context().Value("username").(string)

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	log.Printf("workspaceID: %s; formID: %s", workspaceID, formID)

	submissionData := []SubmissionData{}
	conf := config.Get()
	var err error
	err = conf.DB.View(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		c := formBucket.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			f := SubmissionData{}

			err = json.Unmarshal(v, &f)
			if err != nil {
				log.Println(err)
			}
			submissionData = append(submissionData, f)
		}

		return nil

	})
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}

}

func GetSubmissionInfoHandler(w http.ResponseWriter, r *http.Request) {
	// username := r.Context().Value("username").(string)
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")
	// log.Println(workspaceID)

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}
	log.Println(r.URL.String())

	submissionData := SubmissionData{}
	conf := config.Get()
	err = conf.DB.View(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		err = json.Unmarshal(formBucket.Get(itob(submissionID)), &submissionData)
		if err != nil {
			log.Println(err)
		}

		return nil

	})
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}

}
