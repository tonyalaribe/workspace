package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func UsersAndWorkspaceRoles(w http.ResponseWriter, r *http.Request) {
	db := config.Get().Database
	users, err := db.GetAllUsers()
	if err != nil {
		log.Println(err)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err)
	}
}

func ChangeUserWorkspacePermission(w http.ResponseWriter, r *http.Request) {
	db := config.Get().Database

	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	permissions := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&permissions)
	if err != nil {
		log.Println(err)
	}

	user, err := db.GetUserByEmail(permissions["email"].(string))
	if err != nil {
		log.Println(err)
	}
	role := workspaceID + "-" + permissions["role"].(string)
	user.Roles = append(user.Roles, role)
	err = db.CreateUser(user)
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

func SetupSuperAdmin(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()
	db := conf.Database
	adminUsername := r.URL.Query().Get("u")
	adminUser, err := db.GetUser(adminUsername)
	if err != nil {
		log.Println(err)
	}

	roles := []string{"superadmin"}
	patchObject := gabs.New()
	patchObject.SetP(roles, "app_metadata.roles")

	byteReader := bytes.NewReader(patchObject.Bytes())

	req, err := http.NewRequest("PATCH", "https://emikra.auth0.com/api/v2/users/"+adminUser.ProviderUserID, byteReader)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.Auth0ApiToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	jsonDecoder := json.NewDecoder(res.Body)

	responseObject, err := gabs.ParseJSONDecoder(jsonDecoder)
	if err != nil {
		log.Println(err)
	}

	user := database.User{}
	user.Username = responseObject.Path("username").Data().(string)
	user.Email = responseObject.Path("email").Data().(string)
	user.Name = responseObject.Path("name").Data().(string)
	user.ProviderUserID = responseObject.Path("user_id").Data().(string)

	rolesInterface := responseObject.Path("app_metadata.roles").Data().([]interface{})
	for _, v := range rolesInterface {
		user.Roles = append(user.Roles, v.(string))
	}

	err = db.CreateUser(user)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte("success "))
}
