package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

//CreateUser Creates a user in the database
func (boltDBProvider *BoltDBProvider) CreateUser(user database.User) error {
	userByte, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(boltDBProvider.UsersBucket))
		err = usersBucket.Put([]byte(user.Username), userByte)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

//GetUser returns a user given the username
func (boltDBProvider *BoltDBProvider) GetUser(username string) (database.User, error) {

	var userByte []byte
	existingUser := database.User{}

	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(boltDBProvider.UsersBucket))
		userByte = usersBucket.Get([]byte(username))
		return nil
	})

	if err != nil {
		return existingUser, err
	}

	err = json.Unmarshal(userByte, &existingUser)
	if err != nil {
		return existingUser, err
	}
	return existingUser, nil
}

//GetsUserByEmail returns a user givenits username
func (boltDBProvider *BoltDBProvider) GetUserByEmail(email string) (database.User, error) {

	var userByte []byte
	existingUser := database.User{}

	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(boltDBProvider.UsersBucket))
		c := usersBucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			currUser := database.User{}
			err := json.Unmarshal(v, &currUser)
			if err != nil {
				log.Println(err)
			}
			if currUser.Email == email {
				existingUser = currUser
			}

		}
		return nil
	})

	if err != nil {
		return existingUser, err
	}

	err = json.Unmarshal(userByte, &existingUser)
	return existingUser, err
}

//GetAllUsers returns all users from the datababse
func (boltDBProvider *BoltDBProvider) GetAllUsers() ([]database.User, error) {
	users := []database.User{}

	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		usersBucket := tx.Bucket([]byte(boltDBProvider.UsersBucket))
		err := usersBucket.ForEach(func(k []byte, v []byte) error {
			user := database.User{}
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
