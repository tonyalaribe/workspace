package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/filePersistence"
)

type Config struct {
	RootDirectory     string
	BoltFile          string
	SubmissionsBucket []byte
	DB                *bolt.DB
	FileManager       FileManager
	RolesManager      *gorbac.RBAC
}

type FileManager interface {
	Save(name string, workspace string, b64Data string) (string, error)
}

var (
	config Config
)

const (
	FORMS_METADATA       = "forms_metadata"
	WORKSPACES_METADATA  = "workspaces_metadata"
	WORKSPACES_CONTAINER = "workspaces_container"
	USERS_BUCKET         = "users_bucket"
)

//Using Init not init, so i can manually determine when the content of config are initalized, as opposed to initializing whenever the package is imported (initialization should happen at app startup, which is only when imported by the main.go file).
func Init() {
	config.RootDirectory = filepath.Join(".", "data")
	config.BoltFile = filepath.Join(config.RootDirectory, "workspace.db")
	config.SubmissionsBucket = []byte("submissions")

	os.MkdirAll(config.RootDirectory, os.ModePerm)
	db, err := bolt.Open(config.BoltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(WORKSPACES_METADATA))
		tx.CreateBucketIfNotExists([]byte(WORKSPACES_CONTAINER))
		tx.CreateBucketIfNotExists([]byte(USERS_BUCKET))
		return nil
	})

	log.Println(db.GoString())
	config.DB = db

	config.FileManager = filePersistence.FilePersister{RootDirectory: config.RootDirectory}

	config.RolesManager = gorbac.New()

	// superAdmin := gorbac.NewStdRole("superadmin")

	workspace2Admin := gorbac.NewStdRole("workspace-2-admin")
	workspace2Spectator := gorbac.NewStdRole("workspace-2-spectator")
	workspace2Editor := gorbac.NewStdRole("workspace-2-editor")
	workspace2Supervisor := gorbac.NewStdRole("workspace-2-supervisor")

	viewWorkspace2 := gorbac.NewStdPermission("view-workspace-2")
	editWorkspace2 := gorbac.NewStdPermission("edit-workspace-2")
	approveWorkspace2 := gorbac.NewStdPermission("approve-workspace-2")
	changeRolesWorkspace2 := gorbac.NewStdPermission("change-roles-workspace-2")

	workspace2Spectator.Assign(viewWorkspace2)

	workspace2Editor.Assign(viewWorkspace2)
	workspace2Editor.Assign(editWorkspace2)

	workspace2Supervisor.Assign(viewWorkspace2)
	workspace2Supervisor.Assign(editWorkspace2)
	workspace2Supervisor.Assign(approveWorkspace2)

	workspace2Admin.Assign(viewWorkspace2)
	workspace2Admin.Assign(editWorkspace2)
	workspace2Admin.Assign(approveWorkspace2)
	workspace2Admin.Assign(changeRolesWorkspace2)

}

func SavePermissions() {
	// Persist the change
	// map[RoleId]PermissionIds
	jsonOutputRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonOutputInher := make(map[string][]string)
	SaveJsonHandler := func(r gorbac.Role, parents []string) error {
		// WARNING: Don't use gorbac.RBAC instance in the handler,
		// otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.(*gorbac.StdRole).Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID()] = permissions
		jsonOutputInher[r.ID()] = parents
		return nil
	}
	if err := gorbac.Walk(Get().RolesManager, SaveJsonHandler); err != nil {
		log.Fatalln(err)
	}
	//
	// // Save roles information
	// if err := SaveJson("new-roles.json", &jsonOutputRoles); err != nil {
	// 	log.Fatal(err)
	// }
	// // Save inheritance information
	// if err := SaveJson("new-inher.json", &jsonOutputInher); err != nil {
	// 	log.Fatal(err)
	// }

}

func SaveJSON() {

}

func Get() *Config {
	return &config
}
