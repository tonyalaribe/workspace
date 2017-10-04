package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/actions"
	"gitlab.com/middlefront/workspace/database"
)

type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

//CreateWorkspaceHandler create workspace with database.WorkSpace{} as body.
func CreateWorkspaceHandler(w http.ResponseWriter, r *http.Request) {

	workspaceData := database.WorkSpace{}
	user := r.Context().Value("user").(database.User)

	err := json.NewDecoder(r.Body).Decode(&workspaceData)
	if err != nil {
		log.Println(err)
	}

	message := make(map[string]interface{})
	err = actions.CreateWorkspace(workspaceData, user)
	if err != nil {
		log.Println(err)
		message["code"] = http.StatusInternalServerError
		message["message"] = err.Error()
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	message["code"] = 200
	message["message"] = "success"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Println(err)
	}
}

//GetWorkspacesHandler Get workspaces a user has access to
func GetWorkspacesHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(database.User)

	workspaces, err := actions.GetWorkspaces(user)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(workspaces)
	if err != nil {
		log.Println(err)
	}
}

//GetWorkspaceUsersAndRolesHandler get users and their roles attached to a workspace
func GetWorkspaceUsersAndRolesHandler(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("w")

	users, err := actions.GetWorkspaceUsersAndRoles(workspaceID)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err)
	}
}

//GetWorkspaceBySlugHandler Get a workspaces details given the workspaceID(slug) as param.
func GetWorkspaceBySlugHandler(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")

	workspace, err := actions.GetWorkspaceBySlug(workspaceID)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(workspace)
	if err != nil {
		log.Println(err)
	}
}
