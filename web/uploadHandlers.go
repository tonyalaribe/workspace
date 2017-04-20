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

	submission := SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", submission)

	conf := config.Get()

	var workspaceInfoByte []byte
	conf.DB.Update(func(tx *bolt.Tx) error {
		workspacesBucket, err := tx.CreateBucketIfNotExists([]byte(config.WORKSPACES_BUCKET))
		if err != nil {
			log.Println(err)
		}
		workspaceInfoByte = workspacesBucket.Get([]byte(workspaceID))

		return nil
	})

	workspaceInfo := make(map[string]interface{})
	err = json.Unmarshal(workspaceInfoByte, &workspaceInfo)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", workspaceInfo)

	schema := workspaceInfo["jsonschema"].(map[string]interface{})

	for k, v := range submission.FormData {
		log.Printf("key: %+v, value: %+v", k, "v")
		schemaObject := schema["properties"].(map[string]interface{})[k].(map[string]interface{})
		log.Println(schemaObject)
		switch schemaObject["type"].(string) {
		case "string":
			log.Println("processing a string type")
			log.Println(schemaObject["format"].(string))
			switch schemaObject["format"].(string) {
			case "data-uri", "data-url":
				//file formatting
				log.Println("data-uri processing")
				pathToSubmission := filepath.Join(conf.RootDirectory, username, submission.SubmissionName)
				os.MkdirAll(pathToSubmission, os.ModePerm)
				fullPath := Base64ToFileSystem(v.(string), pathToSubmission)
				submission.FormData[k] = fullPath
				break
			default:
				submission.FormData[k] = v.(string)
			}

		case "integer":
			//Using type float64 due to compiler complaints when handling integer types
			submission.FormData[k] = submission.FormData[k].(float64)
		default:
			log.Printf("%+v", schemaObject)
		}
	}

	log.Println(submission)

	// formData.CreatedBy = username
	// formData.CreationDate = time.Now().Format(time.RFC1123)
	//
	//
	// for i, file := range formData.Files {
	// 	pathToUser := filepath.Join(conf.RootDirectory, username, formData.SubmissionName)
	// 	os.MkdirAll(pathToUser, os.ModePerm)
	// 	filepath := filepath.Join(pathToUser, file.Name)
	// 	Base64ToFileSystem(file.File, filepath)
	//
	// 	formData.Files[i].File = ""
	// 	formData.Files[i].Path = filepath
	// 	formData.Files[i].CreatedBy = username
	// 	formData.Files[i].UploadDate = time.Now().Format(time.RFC1123)
	//
	// }

	/*Save to boltdb*/
	err = conf.DB.Update(func(tx *bolt.Tx) error {

		workspaceBucket, err := tx.CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}

		userBucket, err := workspaceBucket.CreateBucketIfNotExists([]byte(username))
		if err != nil {
			log.Println(err)
		}

		nextID, err := userBucket.NextSequence()
		if err != nil {
			log.Println(err)
		}

		submission.ID = int(nextID)
		dataByte, err := json.Marshal(submission)
		if err != nil {
			log.Println(err)
		}

		err = userBucket.Put(itob(int(nextID)), dataByte)
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

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}
	log.Println(submissionID)

	submissionData := SubmissionData{}

	conf := config.Get()

	err = conf.DB.View(func(tx *bolt.Tx) error {
		workspacesBucket := tx.Bucket([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}

		b := workspacesBucket.Bucket([]byte(username))

		err = json.Unmarshal(b.Get(itob(submissionID)), &submissionData)
		if err != nil {
			log.Println(err)
		}

		return nil

	})
	if err != nil {
		log.Println(err)
	}

	formData := SubmissionData{}
	err = json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}

	submissionData.Status = formData.Status
	submissionData.LastModified = formData.LastModified
	submissionData.FormData = formData.FormData

	// for i, file := range formData.Files {
	// 	pathToUser := filepath.Join(conf.RootDirectory, username, formData.SubmissionName)
	// 	os.MkdirAll(pathToUser, os.ModePerm)
	// 	filepath := filepath.Join(pathToUser, file.Name)
	// 	Base64ToFileSystem(file.File, filepath)
	//
	// 	oldFile := formData.Files[i]
	// 	oldFile.File = ""
	// 	oldFile.Path = filepath
	// 	oldFile.CreatedBy = username
	// 	oldFile.UploadDate = time.Now().Format(time.RFC1123)
	// 	submissionData.Files = append(submissionData.Files, oldFile)

	// }
	//
	/*Save to boltdb*/

	err = conf.DB.Update(func(tx *bolt.Tx) error {
		workspacesBucket, err := tx.CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}

		bucket, err := workspacesBucket.CreateBucketIfNotExists([]byte(username))
		if err != nil {
			log.Println(err)
		}

		dataByte, err := json.Marshal(submissionData)
		if err != nil {
			log.Println(err)
		}

		err = bucket.Put(itob(submissionID), dataByte)
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

func GetMySubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	log.Println(workspaceID)

	submissionData := []SubmissionData{}

	conf := config.Get()

	var err error
	err = conf.DB.View(func(tx *bolt.Tx) error {

		workspaceBucket := tx.Bucket([]byte(workspaceID))

		b := workspaceBucket.Bucket([]byte(username))

		c := b.Cursor()

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
	username := r.Context().Value("username").(string)
	httprouterParams := r.Context().Value("params").(httprouter.Params)

	workspaceID := httprouterParams.ByName("workspaceID")
	// log.Println(workspaceID)

	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}

	log.Println(r.URL.String())

	submissionData := SubmissionData{}

	conf := config.Get()

	err = conf.DB.View(func(tx *bolt.Tx) error {
		workspaceBucket := tx.Bucket([]byte(workspaceID))
		b := workspaceBucket.Bucket([]byte(username))

		err = json.Unmarshal(b.Get(itob(submissionID)), &submissionData)
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
