package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FormData struct {
	Files []struct {
		File string `json:"file"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"files"`
	SubmissionName string `json:"submissionName"`
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
	log.Println(username)
	formData := FormData{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", formData)
	for _, file := range formData.Files {
		pathToUser := filepath.Join(".", "data", username)
		os.MkdirAll(pathToUser, os.ModePerm)
		filepath := filepath.Join(pathToUser, formData.SubmissionName+"^^"+file.Name)

		Base64ToFileSystem(file.File, filepath)
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
	log.Println(username)
	pathToUser := filepath.Join(".", "data", username)
	os.MkdirAll(pathToUser, os.ModePerm)

	type File struct {
		SubmissionName string
		FileName       string
		CreatedBy      string
	}
	submissionData := []File{}

	files, _ := ioutil.ReadDir(pathToUser)
	for _, f := range files {

		filename := f.Name()
		splitFilename := strings.Split(filename, "^^")
		fmt.Println(f.Name())
		submissionData = append(submissionData, File{
			SubmissionName: splitFilename[0],
			FileName:       splitFilename[1],
			CreatedBy:      username,
		})
	}
	w.Header().Set("Content-type", "application/json")

	err := json.NewEncoder(w).Encode(submissionData)
	if err != nil {
		log.Println(err)
	}

}
