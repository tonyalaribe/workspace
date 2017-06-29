package web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/middlefront/workspace/config"
)

type User struct {
	ProviderUserID    string
	Username          string
	Name              string
	Email             string
	Roles             []string
	CurrentRoleString string
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

func (user User) GetByEmail(email string) (User, error) {
	conf := config.Get()
	result := User{}

	err := conf.DB.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(config.USERS_BUCKET))

		c := usersBucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			currUser := User{}
			err := json.Unmarshal(v, &currUser)
			if err != nil {
				log.Println(err)
			}
			if currUser.Email == email {
				result = currUser
			}

		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (user User) GetAll() ([]User, error) {
	conf := config.Get()
	users := []User{}

	err := conf.DB.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(config.USERS_BUCKET))
		err := usersBucket.ForEach(func(k []byte, v []byte) error {
			user := User{}
			err := json.Unmarshal(v, &user)
			if err != nil {
				log.Println(err)
			}
			users = append(users, user)

			return nil
		})
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	if err != nil {
		return users, err
	}

	return users, nil
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

func UsersAndWorkspaceRoles(w http.ResponseWriter, r *http.Request) {
	users, err := User{}.GetAll()
	if err != nil {
		log.Println(err)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err)
	}
}

func ChangeUserWorkspacePermission(w http.ResponseWriter, r *http.Request) {
	httprouterParams := r.Context().Value("params").(httprouter.Params)
	workspaceID := httprouterParams.ByName("workspaceID")
	permissions := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&permissions)
	if err != nil {
		log.Println(err)
	}

	user, err := User{}.GetByEmail(permissions["email"].(string))
	if err != nil {
		log.Println(err)
	}
	role := workspaceID + "-" + permissions["role"].(string)
	user.Roles = append(user.Roles, role)
	err = user.Create()
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
	adminUsername := r.URL.Query().Get("u")
	adminUser, err := User{}.Get(adminUsername)
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
	user := User{}

	user.Username = responseObject.Path("username").Data().(string)
	user.Email = responseObject.Path("email").Data().(string)
	user.Name = responseObject.Path("name").Data().(string)
	user.ProviderUserID = responseObject.Path("user_id").Data().(string)

	rolesInterface := responseObject.Path("app_metadata.roles").Data().([]interface{})
	for _, v := range rolesInterface {
		user.Roles = append(user.Roles, v.(string))
	}

	err = user.Create()
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte("success "))
}
