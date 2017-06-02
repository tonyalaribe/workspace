package web

import (
	"encoding/json"
	"log"

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
	log.Println(userByte)
	if err != nil {
		return existingUser, err
	}
	log.Println(userByte)
	err = json.Unmarshal(userByte, &existingUser)
	if err != nil {
		return existingUser, err
	}
	return existingUser, nil
}
