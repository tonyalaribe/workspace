package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/metal3d/go-slugify"
	"github.com/mikespook/gorbac"
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
	user := r.Context().Value("user").(User)

	err := json.NewDecoder(r.Body).Decode(&workspaceData)
	if err != nil {
		log.Println(err)
	}

	workspaceData.Creator = user.Username
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

	spectator := gorbac.NewStdRole(workspaceData.ID + "-spectator")
	spectator.Assign(gorbac.NewStdPermission("view-" + workspaceData.ID))
	conf.RolesManager.Add(spectator)

	editor := gorbac.NewStdRole(workspaceData.ID + "-editor")
	editor.Assign(gorbac.NewStdPermission("edit-" + workspaceData.ID))
	conf.RolesManager.Add(editor)

	supervisor := gorbac.NewStdRole(workspaceData.ID + "-supervisor")
	supervisor.Assign(gorbac.NewStdPermission("approve-" + workspaceData.ID))
	conf.RolesManager.Add(supervisor)

	admin := gorbac.NewStdRole(workspaceData.ID + "-admin")
	admin.Assign(gorbac.NewStdPermission("admin-" + workspaceData.ID))
	conf.RolesManager.Add(admin)

	conf.RolesManager.SetParent(workspaceData.ID+"-editor", workspaceData.ID+"-spectator")
	conf.RolesManager.SetParent(workspaceData.ID+"-supervisor", workspaceData.ID+"-editor")
	conf.RolesManager.SetParent(workspaceData.ID+"-admin", workspaceData.ID+"-supervisor")

	conf.RolesManager.SetParent("superadmin", workspaceData.ID+"-admin")

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
	user := r.Context().Value("user").(User)
	log.Printf("%#v", user)

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

	finalWorkspaces := []WorkSpace{}
	for _, v := range workspaces {
		workspacePermissionString := "view-" + v.ID
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)
		if gorbac.AnyGranted(conf.RolesManager, user.Roles, workspacePermission, nil) {
			finalWorkspaces = append(finalWorkspaces, v)
		}
	}
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(finalWorkspaces)
	if err != nil {
		log.Println(err)
	}
}

func GetWorkspaceUsersAndRolesHandler(w http.ResponseWriter, r *http.Request) {
	// user := r.Context().Value("user").(User)
	// log.Printf("%#v", user)
	workspaceID := r.URL.Query().Get("w")

	workspace := WorkSpace{}
	users := []User{}

	conf := config.Get()
	conf.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(config.USERS_BUCKET))
		b.ForEach(func(_ []byte, v []byte) error {

			user := User{}
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}
			users = append(users, user)

			return nil
		})

		w := tx.Bucket([]byte(config.WORKSPACES_METADATA))
		wByte := w.Get([]byte(workspaceID))
		err := json.Unmarshal(wByte, &workspace)
		if err != nil {
			return err
		}
		return nil
	})
	log.Println(users)
	finalUsers := []User{}
	for _, u := range users {
		log.Println(u)
		workspacePermissionString := "view-" + workspace.ID
		log.Println(workspacePermissionString)
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)

		if gorbac.AnyGranted(conf.RolesManager, u.Roles, workspacePermission, nil) {
			log.Println("granted final users access")
			finalUsers = append(finalUsers, u)
		}
	}

	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(finalUsers)
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
