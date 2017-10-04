package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/actions"
)

//UsersAndWorkspaceRoles Get users and their roles in the workspace
func UsersAndWorkspaceRoles(w http.ResponseWriter, r *http.Request) {
	users, err := actions.GetUsersAndWorkspaceRoles()
	if err != nil {
		log.Println(err)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err)
	}
}

//ChangeUserWorkspacePermission chamges workspace permissions, but only if the user making the action has enough (admin or similar) permissions
func ChangeUserWorkspacePermission(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	permissions := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&permissions)
	if err != nil {
		log.Println(err)
	}

	err = actions.ChangeUserWorkspacePermission(workspaceID, permissions)
	if err != nil {
		log.Println(err)
	}

	response := map[string]string{}
	response["message"] = "Updated User Roles Successfully"
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}
