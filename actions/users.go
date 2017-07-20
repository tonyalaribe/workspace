package actions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"

	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func GetUsersAndWorkspaceRoles() ([]database.User, error) {
	db := config.Get().Database
	users, err := db.GetAllUsers()
	return users, err
}

func ChangeUserWorkspacePermission(workspaceID string, permissions map[string]interface{}) error {
	db := config.Get().Database
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
	return err
}

func SetupSuperAdmin(adminUsername string) error {
	conf := config.Get()
	db := conf.Database
	log.Println(adminUsername)
	adminUser, err := db.GetUser(adminUsername)
	if err != nil {
		log.Println(err)
		return err
	}

	roles := []string{"superadmin"}
	patchObject := gabs.New()
	patchObject.SetP(roles, "app_metadata.roles")

	byteReader := bytes.NewReader(patchObject.Bytes())

	req, err := http.NewRequest("PATCH", "https://emikra.auth0.com/api/v2/users/"+adminUser.ProviderUserID, byteReader)
	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+conf.Auth0ApiToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	jsonDecoder := json.NewDecoder(res.Body)

	responseObject, err := gabs.ParseJSONDecoder(jsonDecoder)
	if err != nil {
		log.Println(err)
		return err
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
		return err
	}
	return nil
}
