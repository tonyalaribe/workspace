package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/config"
)

type FormData struct {
	Files []struct {
		File string `json:"file"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"files"`
	SubmissionName string `json:"submissionName"`
	Status         string `json:"status"`
}

type File struct {
	SubmissionName string
	FileName       string
	FilePath       string
	CreatedBy      string
	UploadDate     string
	Status         string
}

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

	formData := FormData{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}

	conf := config.Get()

	for _, file := range formData.Files {
		pathToUser := filepath.Join(conf.RootDirectory, username, formData.SubmissionName)
		os.MkdirAll(pathToUser, os.ModePerm)
		filepath := filepath.Join(pathToUser, file.Name)
		Base64ToFileSystem(file.File, filepath)

		data := File{
			FileName:       file.Name,
			FilePath:       filepath,
			CreatedBy:      username,
			UploadDate:     time.Now().Format(time.RFC1123),
			SubmissionName: formData.SubmissionName,
			Status:         formData.Status,
		}
		dataByte, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
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
		err = bucket.Put(itob(int(nextID)), dataByte)
		if err != nil {
			log.Println(err)
		}
		tx.Commit()
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

	submissionData := []File{}

	conf := config.Get()

	var err error
	err = conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(username))

		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			f := File{}

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
