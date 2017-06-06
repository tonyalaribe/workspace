package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/config"
)

type User struct {
	ProviderUserID string
	Username       string
	Name           string
	Email          string
	Roles          []string
}

func SetupSuperAdmin(w http.ResponseWriter, r *http.Request) {
	conf := config.Get()
	adminUsername := r.URL.Query().Get("u")
	adminUser, err := User{}.Get(adminUsername)
	if err != nil {
		log.Println(err)
	}

	roles := []string{"superadmin"}
	patchObject := gabs.New()
	patchObject.SetP(roles, "app_metadata.roles")

	byteReader := bytes.NewReader(patchObject.Bytes())

	log.Println(adminUser.ProviderUserID)
	log.Println(adminUser)
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
	user := User{}

	// log.Println(responseObject.String())
	user.Username = responseObject.Path("username").Data().(string)
	user.Email = responseObject.Path("email").Data().(string)
	user.Name = responseObject.Path("name").Data().(string)
	user.ProviderUserID = responseObject.Path("user_id").Data().(string)

	rolesInterface := responseObject.Path("app_metadata.roles").Data().([]interface{})
	for _, v := range rolesInterface {
		user.Roles = append(user.Roles, v.(string))
	}

	log.Printf("%#v", user)
	err = user.Create()
	if err != nil {
		log.Println(err)
	}
}

func (user User) Create() error {
	conf := config.Get()
	userByte, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = conf.DB.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(config.USERS_BUCKET))
		err := usersBucket.Put([]byte(user.Username), userByte)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (user User) Get(username string) (User, error) {
	conf := config.Get()
	var userByte []byte
	existingUser := User{}

	err := conf.DB.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(config.USERS_BUCKET))
		userByte = usersBucket.Get([]byte(username))
		return nil
	})

	if err != nil {
		return existingUser, err
	}

	err = json.Unmarshal(userByte, &existingUser)
	log.Println(existingUser)
	if err != nil {
		return existingUser, err
	}
	return existingUser, nil
}
