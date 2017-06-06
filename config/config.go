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
	Auth0ApiToken     string
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

	superAdmin := gorbac.NewStdRole("superadmin")

	workspace2Admin := gorbac.NewStdRole("workspace-2-admin")
	workspace2Spectator := gorbac.NewStdRole("workspace-2-spectator")
	workspace2Editor := gorbac.NewStdRole("workspace-2-editor")
	workspace2Supervisor := gorbac.NewStdRole("workspace-2-supervisor")

	superadminRole := gorbac.NewStdPermission("super-admin")

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

	superAdmin.Assign(superadminRole)

	config.RolesManager.Add(workspace2Admin)
	config.RolesManager.Add(workspace2Spectator)
	config.RolesManager.Add(workspace2Editor)
	config.RolesManager.Add(workspace2Supervisor)
	config.RolesManager.Add(superAdmin)

	config.Auth0ApiToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlF6TTBOVFl4TXpVek9VRXpNa1ExUlVNNU5UQTFSRGN3TURaQlF6WTRPVEkwUWpreFJUVTNSZyJ9.eyJpc3MiOiJodHRwczovL2VtaWtyYS5hdXRoMC5jb20vIiwic3ViIjoiWlc5WTJWSlh3TWRqZkliT3lVT1QzZ0paTk5KZm9PQk9AY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vZW1pa3JhLmF1dGgwLmNvbS9hcGkvdjIvIiwiZXhwIjoxNTgyODc5ODM5LCJpYXQiOjE0OTY0Nzk4MzksInNjb3BlIjoicmVhZDpjbGllbnRfZ3JhbnRzIGNyZWF0ZTpjbGllbnRfZ3JhbnRzIGRlbGV0ZTpjbGllbnRfZ3JhbnRzIHVwZGF0ZTpjbGllbnRfZ3JhbnRzIHJlYWQ6dXNlcnMgdXBkYXRlOnVzZXJzIGRlbGV0ZTp1c2VycyBjcmVhdGU6dXNlcnMgcmVhZDp1c2Vyc19hcHBfbWV0YWRhdGEgdXBkYXRlOnVzZXJzX2FwcF9tZXRhZGF0YSBkZWxldGU6dXNlcnNfYXBwX21ldGFkYXRhIGNyZWF0ZTp1c2Vyc19hcHBfbWV0YWRhdGEgY3JlYXRlOnVzZXJfdGlja2V0cyByZWFkOmNsaWVudHMgdXBkYXRlOmNsaWVudHMgZGVsZXRlOmNsaWVudHMgY3JlYXRlOmNsaWVudHMgcmVhZDpjbGllbnRfa2V5cyB1cGRhdGU6Y2xpZW50X2tleXMgZGVsZXRlOmNsaWVudF9rZXlzIGNyZWF0ZTpjbGllbnRfa2V5cyByZWFkOmNvbm5lY3Rpb25zIHVwZGF0ZTpjb25uZWN0aW9ucyBkZWxldGU6Y29ubmVjdGlvbnMgY3JlYXRlOmNvbm5lY3Rpb25zIHJlYWQ6cmVzb3VyY2Vfc2VydmVycyB1cGRhdGU6cmVzb3VyY2Vfc2VydmVycyBkZWxldGU6cmVzb3VyY2Vfc2VydmVycyBjcmVhdGU6cmVzb3VyY2Vfc2VydmVycyByZWFkOmRldmljZV9jcmVkZW50aWFscyB1cGRhdGU6ZGV2aWNlX2NyZWRlbnRpYWxzIGRlbGV0ZTpkZXZpY2VfY3JlZGVudGlhbHMgY3JlYXRlOmRldmljZV9jcmVkZW50aWFscyByZWFkOnJ1bGVzIHVwZGF0ZTpydWxlcyBkZWxldGU6cnVsZXMgY3JlYXRlOnJ1bGVzIHJlYWQ6ZW1haWxfcHJvdmlkZXIgdXBkYXRlOmVtYWlsX3Byb3ZpZGVyIGRlbGV0ZTplbWFpbF9wcm92aWRlciBjcmVhdGU6ZW1haWxfcHJvdmlkZXIgYmxhY2tsaXN0OnRva2VucyByZWFkOnN0YXRzIHJlYWQ6dGVuYW50X3NldHRpbmdzIHVwZGF0ZTp0ZW5hbnRfc2V0dGluZ3MgcmVhZDpsb2dzIHJlYWQ6c2hpZWxkcyBjcmVhdGU6c2hpZWxkcyBkZWxldGU6c2hpZWxkcyB1cGRhdGU6dHJpZ2dlcnMgcmVhZDp0cmlnZ2VycyByZWFkOmdyYW50cyBkZWxldGU6Z3JhbnRzIHJlYWQ6Z3VhcmRpYW5fZmFjdG9ycyB1cGRhdGU6Z3VhcmRpYW5fZmFjdG9ycyByZWFkOmd1YXJkaWFuX2Vucm9sbG1lbnRzIGRlbGV0ZTpndWFyZGlhbl9lbnJvbGxtZW50cyBjcmVhdGU6Z3VhcmRpYW5fZW5yb2xsbWVudF90aWNrZXRzIHJlYWQ6dXNlcl9pZHBfdG9rZW5zIn0.g0BhfrJN6dnQCaw6i1Do-OqlSZBKHOmEfue1Sy0xleKDlvtLujjyy19a1XUcp4IflpRRorhD4D6fTIwc2eeFSJiSbZvFeHW574eAYo88-6y05n_NVpmQVtS1VWJVkaicdd7DhQSWqbBT8VU0UGeB_KjJb5x7xU-yR4V1o6w7kwGZwV0iujJN3XVNuQWwd_t2XqNA-491KNLTPnY00rp_cct8Ru4rZYYNGMskBmGqYyHs9g-XtfcFa0_IeXjUnjaP8SFGK_FjsRER3QWDBaE2CfMDPxxdggmBaiG0yjPmnUbxQLnTpmJKKNTUF-4GM4OzFRZlrQGRa88tWqkXJJnHzQ"
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
