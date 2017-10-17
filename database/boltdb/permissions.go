package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

//SaveRoles persists the current permissions and roles tree
func (boltDBProvider *BoltDBProvider) SaveRoles(roles interface{}) error {

	tx, err := boltDBProvider.db.Begin(true)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	AppMetadataBucket, err := tx.CreateBucketIfNotExists([]byte(boltDBProvider.AppMetadata))
	if err != nil {
		log.Println(err)
	}

	dataByte, err := json.Marshal(roles)
	if err != nil {
		log.Println(err)
	}

	err = AppMetadataBucket.Put([]byte("ROLES_JSON"), dataByte)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

	return nil
}

//GetRoles gets the current roles json string from database
func (boltDBProvider *BoltDBProvider) GetRoles() (string, error) {

	var dataByte []byte
	boltDBProvider.db.View(func(tx *bolt.Tx) error {

		AppMetadataBucket := tx.Bucket([]byte(boltDBProvider.AppMetadata))

		dataByte = AppMetadataBucket.Get([]byte("ROLES_JSON"))

		return nil
	})
	return string(dataByte), nil
}

//SaveInheritance persists tuhe current inheritance tree to database
func (boltDBProvider *BoltDBProvider) SaveInheritance(roles interface{}) error {

	tx, err := boltDBProvider.db.Begin(true)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	AppMetadataBucket, err := tx.CreateBucketIfNotExists([]byte(boltDBProvider.AppMetadata))
	if err != nil {
		log.Println(err)
	}

	dataByte, err := json.Marshal(roles)
	if err != nil {
		log.Println(err)
	}

	err = AppMetadataBucket.Put([]byte("INHERITANCE_JSON"), dataByte)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

	return nil
}

//GetInheritance Get the current Inheritance tree from database
func (boltDBProvider *BoltDBProvider) GetInheritance() (string, error) {

	var dataByte []byte
	boltDBProvider.db.View(func(tx *bolt.Tx) error {

		AppMetadataBucket := tx.Bucket([]byte(boltDBProvider.AppMetadata))

		dataByte = AppMetadataBucket.Get([]byte("INHERITANCE_JSON"))

		return nil
	})

	return string(dataByte), nil
}
