package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/metal3d/go-slugify"
	"gitlab.com/middlefront/workspace/config"
)

type Form struct {
	Creator    string                 `json:"creator"`
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	JSONSchema map[string]interface{} `json:"jsonschema"`
	UISchema   map[string]interface{} `json:"uischema"`
}

func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	formData := Form{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Println(err)
	}

	formData.Creator = r.Context().Value("username").(string)
	formData.ID = slugify.Marshal(formData.Name, true)

	conf := config.Get()
	tx, err := conf.DB.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	currentWorkspaceBucket := tx.Bucket([]byte(conf.WorkspacesContainer)).Bucket([]byte(workspaceID))

	formsMetaDataBucket, err := currentWorkspaceBucket.CreateBucketIfNotExists([]byte(conf.FormsMetadata))
	if err != nil {
		log.Println(err)
	}

	_, err = currentWorkspaceBucket.CreateBucketIfNotExists([]byte(formData.ID))
	if err != nil {
		log.Println(err)
	}

	dataByte, err := json.Marshal(formData)
	if err != nil {
		log.Println(err)
	}

	err = formsMetaDataBucket.Put([]byte(formData.ID), dataByte)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

	message := make(map[string]interface{})
	message["code"] = 200
	message["message"] = "success"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}

func GetFormsHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	forms := []Form{}

	conf := config.Get()
	conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(conf.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(conf.FormsMetadata))
		b.ForEach(func(_ []byte, v []byte) error {

			form := Form{}
			err := json.Unmarshal(v, &form)
			if err != nil {
				return err
			}
			forms = append(forms, form)

			return nil
		})
		return nil
	})

	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(forms)
	if err != nil {
		log.Println(err)
	}
}

func GetFormBySlugHandler(w http.ResponseWriter, r *http.Request) {

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	formID := httprouterParams.ByName("formID")

	formByte := []byte{}

	conf := config.Get()
	conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(conf.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(conf.FormsMetadata))
		formByte = b.Get([]byte(formID))
		return nil
	})

	form := Form{}
	err := json.Unmarshal(formByte, &form)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(form)
	if err != nil {
		log.Println(err)
	}
}
