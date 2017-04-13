package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

const JSONSchema = `{
  "type": "object",
  "required": [
    "firstName",
    "lastName"
  ],
  "properties": {
    "firstName": {
      "type": "string",
      "title": "First name"
    },
    "lastName": {
      "type": "string",
      "title": "Last name"
    },
    "age": {
      "type": "integer",
      "title": "Age"
    },
    "bio": {
      "type": "string",
      "title": "Bio"
    },
    "password": {
      "type": "string",
      "title": "Password",
      "minLength": 3
    }
  }
}`

func Base64ToFileSystem(b64 string, filepath string) {
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile(filepath, byt, 0644)
	if err != nil {
		log.Println(err)
	}
}

func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	submission := SubmissionData{}
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		log.Println(err)
	}

	conf := config.Get()

	schema := make(map[string]interface{})
	err = json.Unmarshal([]byte(JSONSchema), &schema)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", schema)

	for k, v := range submission.FormData {
		schemaObject := schema["properties"].(map[string]interface{})[k].(map[string]interface{})
		switch schemaObject["type"] {
		case "string":
			switch schemaObject["format"] {
			case "data-uri":
				//file formatting
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

	tx, err := conf.DB.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(username))
	if err != nil {
		log.Println(err)
	}

	nextID, err := bucket.NextSequence()
	if err != nil {
		log.Println(err)
	}

	submission.ID = int(nextID)
	dataByte, err := json.Marshal(submission)
	if err != nil {
		log.Println(err)
	}

	err = bucket.Put(itob(int(nextID)), dataByte)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

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
	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}

	submissionData := SubmissionData{}

	conf := config.Get()

	err = conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(username))

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

	tx, err := conf.DB.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(username))
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
	tx.Commit()

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

	submissionData := []SubmissionData{}

	conf := config.Get()

	var err error
	err = conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(username))

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
	submissionID, err := strconv.Atoi(httprouterParams.ByName("submissionID"))
	if err != nil {
		log.Println(err)
	}

	submissionData := SubmissionData{}

	conf := config.Get()

	err = conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(username))

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
