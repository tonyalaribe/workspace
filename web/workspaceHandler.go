package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/metal3d/go-slugify"
	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

func CreateWorkspaceHandler(w http.ResponseWriter, r *http.Request) {

	workspaceData := database.WorkSpace{}
	user := r.Context().Value("user").(database.User)

	err := json.NewDecoder(r.Body).Decode(&workspaceData)
	if err != nil {
		log.Println(err)
	}

	workspaceData.Creator = user.Username
	workspaceData.ID = slugify.Marshal(workspaceData.Name, true)
	workspaceData.Created = int(time.Now().UnixNano() / 1000000) //Get the time since epoch in milli seconds (javascript date compatible)
	conf := config.Get()
	////// Persist workspace

	err = conf.Database.CreateWorkspace(workspaceData)
	if err != nil {
		log.Println(err)
	}

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
	user := r.Context().Value("user").(database.User)
	conf := config.Get()
	//Get Workspaces
	workspaces, err := conf.Database.GetWorkspaces()
	if err != nil {
		log.Println(err)
	}

	finalWorkspaces := []database.WorkSpace{}
	for _, v := range workspaces {
		workspacePermissionString := "view-" + v.ID
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)
		if gorbac.AnyGranted(conf.RolesManager, user.Roles, workspacePermission, nil) {
			finalWorkspaces = append(finalWorkspaces, v)
		}
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(finalWorkspaces)
	if err != nil {
		log.Println(err)
	}
}

func GetWorkspaceUsersAndRolesHandler(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("w")

	conf := config.Get()
	workspace, users, err := conf.Database.GetWorkspaceUsersAndRoles(workspaceID)

	finalUsers := []database.User{}
	for _, u := range users {
		workspacePermissionString := "view-" + workspace.ID
		log.Println(workspacePermissionString)
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)
		for _, v := range u.Roles {
			if conf.RolesManager.IsGranted(v, workspacePermission, nil) {
				u.CurrentRoleString = v
				finalUsers = append(finalUsers, u)
				continue
			}
		}
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(finalUsers)
	if err != nil {
		log.Println(err)
	}
}

func GetWorkspaceBySlugHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	conf := config.Get()
	workspace, err := conf.Database.GetWorkspaceBySlug(workspaceID)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(workspace)
	if err != nil {
		log.Println(err)
	}
}
