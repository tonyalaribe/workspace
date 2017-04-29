package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/metal3d/go-slugify"
	"gitlab.com/middlefront/workspace/config"
)

type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

func CreateWorkspaceHandler(w http.ResponseWriter, r *http.Request) {

	workspaceData := WorkSpace{}

	err := json.NewDecoder(r.Body).Decode(&workspaceData)
	if err != nil {
		log.Println(err)
	}

	workspaceData.Creator = r.Context().Value("username").(string)
	workspaceData.ID = slugify.Marshal(workspaceData.Name, true)
	workspaceData.Created = int(time.Now().UnixNano() / 1000000) //Get the time since epoch in milli seconds (javascript date compatible)

	conf := config.Get()
	tx, err := conf.DB.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	//Create the bucket where forms under this workspace would be stored.
	individualWorkspace, err := tx.Bucket([]byte(config.WORKSPACES_CONTAINER)).CreateBucketIfNotExists([]byte(workspaceData.ID))
	if err != nil {
		log.Println(err)
	}
	_, err = individualWorkspace.CreateBucketIfNotExists([]byte(config.FORMS_METADATA))
	if err != nil {
		log.Println(err)
	}

	metadata_bucket, err := tx.CreateBucketIfNotExists([]byte(config.WORKSPACES_METADATA))
	if err != nil {
		log.Println(err)
	}

	dataByte, err := json.Marshal(workspaceData)
	if err != nil {
		log.Println(err)
	}

	err = metadata_bucket.Put([]byte(workspaceData.ID), dataByte)
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

func GetWorkspacesHandler(w http.ResponseWriter, r *http.Request) {

	workspaces := []WorkSpace{}

	conf := config.Get()
	conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.WORKSPACES_METADATA))
		b.ForEach(func(_ []byte, v []byte) error {

			workspace := WorkSpace{}
			err := json.Unmarshal(v, &workspace)
			if err != nil {
				return err
			}
			workspaces = append(workspaces, workspace)

			return nil
		})
		return nil
	})

	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(workspaces)
	if err != nil {
		log.Println(err)
	}
}

func GetWorkspaceBySlugHandler(w http.ResponseWriter, r *http.Request) {

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	log.Println(workspaceID)
	workspaceByte := []byte{}

	conf := config.Get()
	conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.WORKSPACES_METADATA))
		workspaceByte = b.Get([]byte(workspaceID))
		return nil
	})

	workspace := WorkSpace{}
	err := json.Unmarshal(workspaceByte, &workspace)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(workspace)
	if err != nil {
		log.Println(err)
	}
}
