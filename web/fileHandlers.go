package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func Base64ToFileSystem(b64 string, filename string) {
	byt, err := base64.StdEncoding.DecodeString(strings.Split(b64, "base64,")[1])
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile("./data/"+filename, byt, 0644)
	if err != nil {
		log.Println(err)
	}
}

func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	formData := FormData{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", formData)
	for _, file := range formData.Files {
		filename := file.Name
		Base64ToFileSystem(file.File, filename)
	}

	response := map[string]string{}
	response["message"] = "Upload Success"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}
