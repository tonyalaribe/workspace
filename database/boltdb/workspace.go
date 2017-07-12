package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"

	"gitlab.com/middlefront/workspace/database"
)

func (boltDBProvider *BoltDBProvider) CreateWorkspace(workspaceData database.WorkSpace) error {

	tx, err := boltDBProvider.db.Begin(true)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	//Create the bucket where forms under this workspace would be stored.
	individualWorkspace, err := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).CreateBucketIfNotExists([]byte(workspaceData.ID))
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = individualWorkspace.CreateBucketIfNotExists([]byte(boltDBProvider.FormsMetadata))
	if err != nil {
		log.Println(err)
		return err
	}

	metadata_bucket, err := tx.CreateBucketIfNotExists([]byte(boltDBProvider.WorkspacesMetadata))
	if err != nil {
		log.Println(err)
		return err
	}

	dataByte, err := json.Marshal(workspaceData)
	if err != nil {
		log.Println(err)
		return err
	}

	err = metadata_bucket.Put([]byte(workspaceData.ID), dataByte)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()

	return nil
}

func (boltDBProvider *BoltDBProvider) GetWorkspaces() ([]database.WorkSpace, error) {
	workspaces := []database.WorkSpace{}

	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.WorkspacesMetadata))
		b.ForEach(func(_ []byte, v []byte) error {

			workspace := database.WorkSpace{}
			err := json.Unmarshal(v, &workspace)
			if err != nil {
				return err
			}
			workspaces = append(workspaces, workspace)

			return nil
		})
		return nil
	})

	return workspaces, nil
}

func (boltDBProvider *BoltDBProvider) GetWorkspaceBySlug(workspaceID string) (database.WorkSpace, error) {
	workspaceByte := []byte{}
	workspace := database.WorkSpace{}

	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.WorkspacesMetadata))
		workspaceByte = b.Get([]byte(workspaceID))
		return nil
	})

	err := json.Unmarshal(workspaceByte, &workspace)
	if err != nil {
		log.Println(err)
		return workspace, err
	}

	return workspace, nil
}

func (boltDBProvider *BoltDBProvider) GetWorkspaceUsersAndRoles(workspaceID string) (database.WorkSpace, []database.User, error) {
	workspace := database.WorkSpace{}
	users := []database.User{}

	boltDBProvider.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltDBProvider.UsersBucket))
		b.ForEach(func(_ []byte, v []byte) error {
			user := database.User{}
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})

		w := tx.Bucket([]byte(boltDBProvider.WorkspacesMetadata))
		wByte := w.Get([]byte(workspaceID))
		err := json.Unmarshal(wByte, &workspace)
		if err != nil {
			return err
		}
		return nil
	})

	return workspace, users, nil
}
