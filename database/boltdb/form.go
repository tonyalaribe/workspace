package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"

	"gitlab.com/middlefront/workspace/database"
)

func (boltDBProvider *BoltDBProvider) CreateForm(workspaceID string, formData database.Form) error {

	tx, err := boltDBProvider.db.Begin(true)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	currentWorkspaceBucket := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID))

	formsMetaDataBucket, err := currentWorkspaceBucket.CreateBucketIfNotExists([]byte(boltDBProvider.FormsMetadata))
	if err != nil {
		log.Println(err)
	}

	_, err = currentWorkspaceBucket.CreateBucketIfNotExists([]byte(formData.ID))
	if err != nil {
		log.Println(err)
	}

	dataByte, err := json.Marshal(formData)
	if err != nil {
		log.Println(err)
	}

	err = formsMetaDataBucket.Put([]byte(formData.ID), dataByte)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()

	return nil
}

func (boltDBProvider *BoltDBProvider) GetForms(workspaceID string) ([]database.Form, error) {

	forms := []database.Form{}
	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(boltDBProvider.FormsMetadata))
		b.ForEach(func(_ []byte, v []byte) error {
			form := database.Form{}
			err := json.Unmarshal(v, &form)
			if err != nil {
				return err
			}
			forms = append(forms, form)

			return nil
		})
		return nil
	})
	return forms, nil
}

func (boltDBProvider *BoltDBProvider) GetFormBySlug(workspaceID, formID string) (database.Form, error) {

	form := database.Form{}
	formByte := []byte{}

	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(boltDBProvider.FormsMetadata))
		formByte = b.Get([]byte(formID))
		return nil
	})

	err := json.Unmarshal(formByte, &form)
	if err != nil {
		log.Println(err)
	}

	return form, nil
}

func (boltDBProvider *BoltDBProvider) GetFormJSONBySlug(workspaceID, formID string) ([]byte, error) {

	formByte := []byte{}

	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(boltDBProvider.FormsMetadata))
		formByte = b.Get([]byte(formID))
		return nil
	})

	return formByte, nil
}
